package src

import "bufio"
import "os"

func Parse(path string, job func(line string, line_number int) string) {
	// Input file
	file, err := os.Open(path)
    assert(err)
	defer file.Close()
    var scanner *bufio.Scanner = bufio.NewScanner(file)

	// Output file
    output_file, err := os.Create(path + ".tmp")
	assert(err)
    defer output_file.Close()
	var writer *bufio.Writer = bufio.NewWriter(output_file)
    defer writer.Flush()

	var line string
	var line_number int = 1
    for scanner.Scan() {
        line = scanner.Text()
		line = job(line, line_number)
		line_number++

		// Write line to output_file
		_, err := writer.WriteString(line + "\n")
		assert(err)
    }

	// Detect table if it was at the end of the file
	job("", line_number)

	err = scanner.Err()
	assert(err)

	// move file
}

// - [v] while looping, start writing a new file, .tmp, Q: write per line or per chunk? how big are the chunks?
// - [v] mark table start with processed mark
// - [ ] on table end, call into REST API
// - [ ] if true, continue
// - [ ] if false, delete .tmp file, panic
// - [ ] if no errors, overwrite file with .tmp file
