/*
 * Copyright (C) 2024 Riyyi
 *
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

// go mod init worklog
// go build
// go run .
// go mod tidy

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/alexflint/go-arg"

	"worklog/src"
)

type Args struct {
	Decl string   `arg:"-d,--decl" help:"Generate travel declaration table" placeholder:"MONTH"`
	Process bool  `arg:"-p,--process" help:"Process specified file and call Jira API"`
	File string   `arg:"positional,required" help:"the worklog file to process"`
}

func (Args) Description() string {
	return "worklog - process a worklog file\n"
}

func main() {
	var args Args
	parser := arg.MustParse(&args)

	// File validation
	_, err := os.Stat(args.File);
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
		parser.Fail("file was not readable: " + args.File)
	}

	// Month validation
	var month int
	if args.Decl != "" {
		month, err = strconv.Atoi(args.Decl)
		if err != nil || month < 1 || month > 12 {
			parser.Fail("decl is not a valid month")
		}
	}

	// Execute

	if args.Process {
		var process src.Process = src.MakeProcess()
		var job = func(line string, line_number int) string {
			return process.Process(line, line_number)
		}
		src.File.Parse(args.File, job, true)
	}

	if month > 0 {
		var decl src.Declaration = src.MakeDeclaration(month)
		var job = func(line string, line_number int) string {
			return decl.Generate(line, line_number)
		}
		src.File.Parse(args.File, job, false)
		fmt.Println(decl.Result())
	}
}
