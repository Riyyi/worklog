/*
 * Copyright (C) 2024 Riyyi
 *
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package src

import "bufio"
import "os"

var File file

type file struct {}

func (file) Parse(path string, job func(line string, line_number int) string, overwrite bool) {
	// Input file
	file, err := os.Open(path)
    assert(err)
	defer file.Close()
    var scanner *bufio.Scanner = bufio.NewScanner(file)

	// Output file
	var writer *bufio.Writer
	if overwrite {
		output_file, err := os.Create(path + ".tmp")
		assert(err)
		defer output_file.Close()
		writer = bufio.NewWriter(output_file)
		defer writer.Flush()
	}

	var line string
	var line_number int = 1
    for scanner.Scan() {
        line = scanner.Text()
		line = job(line, line_number)
		line_number++

		// Write line to output_file
		if overwrite && writer != nil {
			_, err := writer.WriteString(line + "\n")
			assert(err)
		}
    }

	// Detect table if it was at the end of the file
	job("", line_number)

	err = scanner.Err()
	assert(err)

	if overwrite {
		err = os.Rename(path + ".tmp", path)
		assert(err)
	}
}
