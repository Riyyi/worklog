package src

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Date struct {
	clock_in string
	location string
	clock_out string

	last_time string
	last_item_id string
	last_description string
}

type Api struct {
	clock_in *regexp.Regexp
	task_line *regexp.Regexp
	clock_out *regexp.Regexp
	dates map[string]Date
}

// Constructor
func MakeApi() Api {
	return Api{
		clock_in: compileRegex(`^\s*\| [0-9]{4}-[0-9]{2}-[0-9]{2} \|\s+IN\s+\|`),
		task_line: compileRegex(`^\s*\| [0-9]{4}-[0-9]{2}-[0-9]{2} \| [0-9]{2}:[0-9]{2} \|`),
		clock_out: compileRegex(`^\s*\| [0-9]{4}-[0-9]{2}-[0-9]{2} \|\s+OUT\s+\|`),
		dates: make(map[string]Date),
	}
}

func (self *Api) Process(line string, line_number int) string {
	var err error
	if self.clock_in.MatchString(line) {
		err = self.parseClockIn(line, line_number)
	} else if self.task_line.MatchString(line) {
		err = self.parseTask(line, line_number)
	} else if self.clock_out.MatchString(line) {
		err = self.parseClockOut(line, line_number)
	}
	assert(err)

	// fmt.Println(line)
	return line
}

// -----------------------------------------

func compileRegex(pattern string) *regexp.Regexp {
    regex, err := regexp.Compile(pattern)
	assert(err)

	return regex
}

func (self *Api) parseLine(line string, line_number int, size int) ([]string, error) {
	var data []string = strings.Split(line, "|")

	if len(data) != size + 2 {
		return nil, errors.New("malformed line " + strconv.Itoa(line_number) + "\n" + line)
	}

	data = data[1:size + 1]
	for i, value := range data {
        data[i] = strings.TrimSpace(value)

		if len(data[i]) == 0 {
			return nil, errors.New("malformed line " + strconv.Itoa(line_number) + "\n" + line)
		}
	}

	return data, nil
}

func (self *Api) parseClockIn(line string, line_number int) error {
	data, err := self.parseLine(line, line_number, 4)
	if err != nil { return err }

	// Set clock_in, location
	var date Date = self.dates[data[0]]
	date.clock_in = data[2]
	date.location = data[3]
	self.dates[data[0]] = date

	return nil
}

func (self *Api) parseTask(line string, line_number int) error {
	data, err := self.parseLine(line, line_number, 5)
	if err != nil { return err }

	var date Date = self.dates[data[0]]

	if date.clock_in == "" {
		return errors.New("no clock-in time found")
	}

	// Call API for the previous task
	if date.last_time != "" && date.last_item_id != "" && date.last_description != "" {
		err = self.callApi(data[0], date.last_time, data[1], date.last_item_id, date.last_description)
	}

	if err != nil { return err }

	// Set last_time, last_item_id, description
	if data[3] == "X" {
		date.last_time = data[1]
		date.last_item_id = data[2]
		date.last_description = data[4]
	} else { // "V", task is already processed
		date.last_time = ""
		date.last_item_id = ""
		date.last_description = ""
	}
	self.dates[data[0]] = date

	return nil
}

func (self *Api) parseClockOut(line string, line_number int) error {
	data, err := self.parseLine(line, line_number, 3)
	if err != nil { return err }

	// Set clock_out
	var date Date = self.dates[data[0]]
	date.clock_out = data[2]
	self.dates[data[0]] = date

	if date.last_time == "" || date.last_item_id == "" || date.last_description == "" {
		return errors.New("no previous task to use clock-out on")
	}

	// Call API for last task of the day
	self.callApi(data[0], date.last_time, date.clock_out, date.last_item_id, date.last_description)

	return nil
}

func (self *Api) callApi(date string, from_time string, to_time string, item_id string, description string) error {
	fmt.Println("API |" + date + "|" + from_time + "|" + to_time + "|" + item_id + "|" + description)

	// parse line
	// call API
	// error checking

	return nil
}

// Example worklog:

// | 2024-07-06 | IN | 08:30 | Office |

// | 2024-07-06 | 09:00 | T1-123 | V | I did nothing! |
// | 2024-07-06 | 09:30 | T1-456 | X | Blabla |
// | 2024-07-06 | 11:00 | T1-789 | X | - |

// | 2024-07-06 | OUT | 13:00 |
