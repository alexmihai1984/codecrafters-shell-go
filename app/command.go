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

type redirect struct {
	redirectType redirectType
	file         string
}

type redirectType int

const (
	Stdout redirectType = iota
	Stderr
	StdoutAppend
	StderrAppend
	None
)

var builtins = map[string]executable{
	"cd":   &cd{},
	"echo": &echo{},
	"exit": &exit{},
	"pwd":  &pwd{},
	"type": &typeCommand{},
}

func parseCmd(cmd string) command {
	parser := buildParser(cmd)
	parts := parser.parse()
	interpreter := buildInterpreter(parts)
	command := interpreter.interpret()

	return command
}

func (c command) execute() {
	stdout := os.Stdout
	stderr := os.Stderr

	for _, redirect := range c.redirects {
		var file *os.File
		var err error
		if redirect.redirectType == StdoutAppend {
			file, err = os.OpenFile(redirect.file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Print("nsh: ", redirect.file, ": No such file or directory\n")
				return
			}
		} else {
			file, err = os.Create(redirect.file)
			if err != nil {
				fmt.Print("nsh: ", redirect.file, ": No such file or directory\n")
				return
			}
		}
		defer file.Close()

		switch redirect.redirectType {
		case Stdout, StdoutAppend:
			stdout = file
		case Stderr:
			stderr = file
		}
	}

	c.executable.setStdout(stdout)
	c.executable.setStderr(stderr)

	c.executable.execute(c.args)
}
