package main

import (
	"fmt"
	"io"
	"os"
)

type command struct {
	verb        string
	args        []string
	stdRedirect string
	errRedirect string
	executable  executable
}

type shell struct {
	stdout io.Writer
	stderr io.Writer
}

var builtins = map[string]executable{
	"cd":   &cd{},
	"echo": &echo{},
	"exit": &exit{},
	"pwd":  &pwd{},
	"type": &typeCommand{},
}

func parseCmd(cmd string) command {
	parser := buildParser(cmd)
	parserResult := parser.parse()

	if parserResult.verb == "" || !parserResult.ok {
		return command{
			"noop",
			nil,
			parserResult.stdRedirect,
			parserResult.errRedirect,
			&noop{},
		}
	}

	executable, ok := builtins[parserResult.verb]
	if !ok {
		return command{
			parserResult.verb,
			parserResult.params,
			parserResult.stdRedirect,
			parserResult.errRedirect,
			&program{
				name: parserResult.verb,
			},
		}
	}

	return command{
		parserResult.verb,
		parserResult.params,
		parserResult.stdRedirect,
		parserResult.errRedirect,
		executable,
	}
}

func (c command) execute() {
	stdout := os.Stdout
	stderr := os.Stderr

	if c.stdRedirect != "" {
		file, err := os.Create(c.stdRedirect)
		if err != nil {
			fmt.Print("nsh: ", c.stdRedirect, ": No such file or directory\n")
			return
		}

		stdout = file
		defer file.Close()
	}

	if c.errRedirect != "" {
		file, err := os.Create(c.errRedirect)
		if err != nil {
			fmt.Print("nsh: ", c.errRedirect, ": No such file or directory\n")
			return
		}

		stderr = file
		defer file.Close()
	}

	c.executable.setStdout(stdout)
	c.executable.setStderr(stderr)

	c.executable.execute(c.args)
}
