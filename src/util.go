/*
 * Copyright (C) 2024 Riyyi
 *
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package src

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var Util util

type util struct {}

func (util) CompileRegex(pattern string) *regexp.Regexp {
    regex, err := regexp.Compile(pattern)
	assert(err)

	return regex
}

func (util) ParseLine(line string, line_number int, size int) ([]string, error) {
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
