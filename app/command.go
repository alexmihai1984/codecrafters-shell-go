package main

import (
	"fmt"
	"os"
)

type command struct {
	verb       string
	args       []string
	redirects  []redirect
	executable executable
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
			parserResult.redirects,
			&noop{},
		}
	}

	executable, ok := builtins[parserResult.verb]
	if !ok {
		return command{
			parserResult.verb,
			parserResult.params,
			parserResult.redirects,
			&program{
				name: parserResult.verb,
			},
		}
	}

	return command{
		parserResult.verb,
		parserResult.params,
		parserResult.redirects,
		executable,
	}
}

func (c command) execute() {
	stdout := os.Stdout
	stderr := os.Stderr

	for _, redirect := range c.redirects {
		file, err := os.Create(redirect.file)
		if err != nil {
			fmt.Print("nsh: ", redirect.file, ": No such file or directory\n")
			return
		}

		switch redirect.redirectType {
		case Stdout:
			stdout = file
		case Stderr:
			stderr = file
		}

		defer file.Close()
	}

	c.executable.setStdout(stdout)
	c.executable.setStderr(stderr)

	c.executable.execute(c.args)
}
