package main

import (
	"fmt"
	"os"
	"os/exec"
)

type program struct {
	name string
}

func (p program) execute(args []string) {
	filepath := findInPath(p.name)
	if filepath == "" {
		fmt.Print(p.name, ": command not found\n")
	}
	cmd := exec.Command(filepath, args...)
	cmd.Args[0] = p.name
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.Run()
}
