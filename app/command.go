package main

import (
	"fmt"
	"strings"
	"unicode"
)

type command struct {
	verb       string
	args       []string
	executable executable
}

var builtins = map[string]executable{
	"cd":   cd{},
	"echo": echo{},
	"exit": exit{},
	"pwd":  pwd{},
	"type": typeCommand{},
}

func parseCmd(cmd string) command {
	parts := processCommand(cmd) //strings.Fields(cmd)

	if len(parts) == 0 {
		return command{
			"noop",
			nil,
			noop{},
		}
	}

	executable, ok := builtins[parts[0]]
	if !ok {
		return command{
			parts[0],
			parts[1:],
			program{
				parts[0],
			},
		}
	}

	return command{
		parts[0],
		parts[1:],
		executable,
	}
}

func processCommand(cmd string) []string {
	var res []string
	var currentWord strings.Builder

	isSingleQuotesEscaped := false
	for _, char := range cmd {
		if char == '\'' {
			isSingleQuotesEscaped = !isSingleQuotesEscaped
			continue
		}

		if unicode.IsSpace(char) {
			if isSingleQuotesEscaped {
				currentWord.WriteRune(char)
				continue
			} else if currentWord.Len() == 0 {
				// ignore
				continue
			} else {
				res = append(res, currentWord.String())
				currentWord.Reset()
			}
		} else {
			currentWord.WriteRune(char)
		}
	}

	if currentWord.Len() > 0 {
		res = append(res, currentWord.String())
	}

	return res
}

func (c command) execute() {
	c.executable.execute(c.args)
}
