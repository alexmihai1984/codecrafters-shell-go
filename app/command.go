package main

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
	parser := buildParser(cmd)
	parts := parser.parse()

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

func (c command) execute() {
	c.executable.execute(c.args)
}
