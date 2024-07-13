/*
 * Copyright (C) 2024 Riyyi
 *
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package src

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
)

type Location int

const (
	NONE Location = iota // 0
	HOME
	OFFICE
	VISIT
)

type Declaration struct {
	month int
	clock_in *regexp.Regexp
	dates [31]Location
	result string
}

// Constructor
func MakeDeclaration(month int) Declaration {
	return Declaration{
		month: month,
		clock_in: Util.CompileRegex(`^\s*\|\s+[0-9]{4}-[0-9]{2}-[0-9]{2}\s+\|\s+IN\s+\|\s+[0-9]{2}:[0-9]{2}\s+\|\s+[a-zA-Z]+\s+\|`),
	}
}

func (self *Declaration) Generate(line string, line_number int) string {
	var err error
	if self.clock_in.MatchString(line) {
		err = self.parseLocation(line, line_number)
	}
	assert(err)

	return line
}

func (self *Declaration) Result() string {
	var result string
	for _, date := range self.dates {
		if date == HOME {
			result += "x\t\t"
		} else if date == OFFICE {
			result +="\tx\t"
		} else if date == VISIT {
			result += "\t\tx"
		}
		result += "\n"
	}
	clipboard.WriteAll(result)
	result = "-home--office--visit-" + result
	result = result + "---------------------"

	return result
}

// -----------------------------------------

func (self *Declaration) parseLocation(line string, line_number int) error {
	data, err := Util.ParseLine(line, line_number, 4)
	if err != nil { return err }

	var month_string string = data[0][5:7]
	month, err := strconv.Atoi(month_string)
	if err != nil || month < 1 || month > 12 {
		return fmt.Errorf("invalid month '%s' on line %d\n%s", month_string, line_number, line)
	}

	var day_string string = data[0][8:]
	day, err := strconv.Atoi(day_string)
	if err != nil || day < 1 || day > 31 {
		return fmt.Errorf("invalid day '%s' on line %d\n%s", day_string, line_number, line)
	}

	if month == self.month {
		var data_month = strings.ToLower(data[3])
		if data_month == strings.ToLower("Home") {
			self.dates[day - 1] = HOME
		} else if data_month == strings.ToLower("Office") {
			self.dates[day - 1] = OFFICE
		} else if data_month == strings.ToLower("Visit") {
			self.dates[day - 1] = VISIT
		} else {
			return fmt.Errorf("invalid location '%s' on line %d\n%s", data[3], line_number, line)
		}
	}

	return nil
}

// | 2024-07-06 | IN | 08:30 | Office |
