package main

import (
	"fmt"
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

	output, err := cmd.Output()
	if err != nil {

	}

	fmt.Print(string(output))
}
