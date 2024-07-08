// go mod init worklog
// go build
// go run .
// go mod tidy

package main

import "errors"
import "os"

import "github.com/alexflint/go-arg"

import "worklog/src"

type Args struct {
	Decl string   `arg:"-d,--decl" help:"Generate travel declaration table" placeholder:"MONTH"`
	Process bool  `arg:"-p,--process" help:"Process specified file and call Jira API"`
	File string   `arg:"positional,required" help:"the worklog file to process"`
}

func (Args) Description() string {
	return "\nworklog - process a worklog file\n"
}

func main() {
	var args Args
	parser := arg.MustParse(&args)

	_, err := os.Stat(args.File);
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
		parser.Fail("file was not readable: " + args.File)
	}

	if args.Process {
		var api src.Api = src.MakeApi()
		var job = func(line string, line_number int) string {
			return api.Process(line, line_number)
		}
		src.Parse(args.File, job)
	}

	if args.Decl != "" {
		// TODO: generate declaration table..
	}
}
