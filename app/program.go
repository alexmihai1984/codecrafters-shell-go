package main

import (
	"fmt"
	"io"
	"os/exec"
)

type program struct {
	shell
	name string
}

func (p *program) setStdout(stdout io.Writer) {
	p.stdout = stdout
}

func (p *program) setStderr(stderr io.Writer) {
	p.stderr = stderr
}

func (p program) execute(args []string) {
	filepath := findInPath(p.name)
	if filepath == "" {
		filepath = findAbsolutePath(p.name)
		if filepath == "" {
			fmt.Fprint(p.stderr, p.name, ": command not found\n")
			return
		}
	}

	cmd := exec.Command(filepath, args...)
	cmd.Args[0] = p.name
	cmd.Stderr = p.stderr
	cmd.Stdout = p.stdout

	cmd.Run()
}
