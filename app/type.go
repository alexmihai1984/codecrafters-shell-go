package main

import (
	"fmt"
	"io"
)

type typeCommand struct {
	shell
}

func (t *typeCommand) setStdout(stdout io.Writer) {
	t.stdout = stdout
}

func (t *typeCommand) setStderr(stderr io.Writer) {
	t.stderr = stderr
}

func (t typeCommand) isBuiltin(cmd string) bool {
	_, ok := builtins[cmd]

	return ok
}

func (t typeCommand) execute(args []string) {
	for _, cmd := range args {
		if t.isBuiltin(cmd) {
			fmt.Fprintln(t.stdout, cmd, "is a shell builtin")
		} else if filepath := findInPath(cmd); filepath != "" {
			fmt.Fprintln(t.stdout, cmd, "is", filepath)
		} else {
			fmt.Fprint(t.stdout, cmd, ": not found\n")
		}
	}
}
