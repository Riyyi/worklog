/*
 * Copyright (C) 2024 Riyyi
 *
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package src

import (
	"errors"
	"regexp"
)

type Date struct {
	clock_in string
	location string
	clock_out string

	last_time string
	last_item_id string
	last_description string
	last_was_processed bool
}

type Process struct {
	clock_in *regexp.Regexp
	task_line *regexp.Regexp
	clock_out *regexp.Regexp
	mark_processed *regexp.Regexp
	dates map[string]Date
}

// Constructor
func MakeProcess() Process {
	return Process{
		clock_in: Util.CompileRegex(`^\s*\|\s+[0-9]{4}-[0-9]{2}-[0-9]{2}\s+\|\s+IN\s+\|`),
		task_line: Util.CompileRegex(`^\s*\|\s+[0-9]{4}-[0-9]{2}-[0-9]{2}\s+\|\s+[0-9]{2}:[0-9]{2}\s+\|`),
		clock_out: Util.CompileRegex(`^\s*\|\s+[0-9]{4}-[0-9]{2}-[0-9]{2}\s+\|\s+OUT\s+\|`),
		mark_processed: Util.CompileRegex(`\|\s+X\s+\|`),
		dates: make(map[string]Date),
	}
}

func (self *Process) Process(line string, line_number int) string {
	var err error
	if self.clock_in.MatchString(line) {
		err = self.parseClockIn(line, line_number)
	} else if self.task_line.MatchString(line) {
		err = self.parseTask(line, line_number)
		// This marks the current task line as "processed",
		// even though we do the API call on the previous task each iteration.
		// When an error occurs, you should take this into account!
		line = self.mark_processed.ReplaceAllString(line, "| V |")
	} else if self.clock_out.MatchString(line) {
		err = self.parseClockOut(line, line_number)
	}
	assert(err)

	return line
}

// -----------------------------------------

func (self *Process) parseClockIn(line string, line_number int) error {
	data, err := Util.ParseLine(line, line_number, 4)
	if err != nil { return err }

	// Set clock_in, location
	var date Date = self.dates[data[0]]
	date.clock_in = data[2]
	date.location = data[3]
	self.dates[data[0]] = date

	return nil
}

func (self *Process) parseTask(line string, line_number int) error {
	data, err := Util.ParseLine(line, line_number, 5)
	if err != nil { return err }

	var date Date = self.dates[data[0]]

	if date.clock_in == "" {
		return errors.New("no clock-in time found")
	}

	// Call API for the previous task
	if !date.last_was_processed && date.last_time != "" && date.last_item_id != "" {
		err = Api.CallApi(data[0], date.last_time, data[1], date.last_item_id, date.last_description)
	}

	if err != nil { return err }

	// Set last_time, last_item_id, description
	if data[3] == "X" {
		date.last_time = data[1]
		date.last_item_id = data[2]
		date.last_description = data[4]
		date.last_was_processed = false
	} else { // "V", task is already processed
		date.last_time = ""
		date.last_item_id = ""
		date.last_description = ""
		date.last_was_processed = true
	}
	self.dates[data[0]] = date

	return nil
}

func (self *Process) parseClockOut(line string, line_number int) error {
	data, err := Util.ParseLine(line, line_number, 3)
	if err != nil { return err }

	// Set clock_out
	var date Date = self.dates[data[0]]
	date.clock_out = data[2]
	self.dates[data[0]] = date

	if !date.last_was_processed && date.last_time == "" || date.last_item_id == "" {
		return errors.New("no previous task to use clock-out on")
	}

	// Call API for last task of the day
	Api.CallApi(data[0], date.last_time, date.clock_out, date.last_item_id, date.last_description)

	return nil
}

// Example worklog:

// | 2024-07-06 | IN | 08:30 | Office |

// | 2024-07-06 | 09:00 | T1-123 | V | I did nothing! |
// | 2024-07-06 | 09:30 | T1-456 | X | Blabla |
// | 2024-07-06 | 11:00 | T1-789 | X | - |

// | 2024-07-06 | OUT | 13:00 |
