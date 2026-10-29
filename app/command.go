package main

import "strings"

type command struct {
	verb       string
	args       []string
	executable executable
}

var executables = map[string]executable{
	"echo": echo{},
	"type": typeCommand{},
	"exit": exit{},
}

func parseCmd(cmd string) command {
	parts := strings.Fields(cmd)

	if len(parts) == 0 {
		return command{
			"noop",
			nil,
			noop{},
		}
	}

	executable, ok := executables[parts[0]]
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

func (c command) execute() {
	c.executable.execute(c.args)
}
