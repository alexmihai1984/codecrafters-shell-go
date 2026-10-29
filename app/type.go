package main

import (
	"fmt"
)

type typeCommand struct{}

var builtins = map[string]bool{
	"echo": true,
	"exit": true,
	"type": true,
}

func isBuiltin(cmd string) bool {
	_, ok := builtins[cmd]

	return ok
}

func (t typeCommand) execute(args []string) {
	for _, cmd := range args {
		if isBuiltin(cmd) {
			fmt.Println(cmd, "is a shell builtin")
		} else if filepath := findInPath(cmd); filepath != "" {
			fmt.Println(cmd, "is", filepath)
		} else {
			fmt.Print(cmd, ": not found\n")
		}
	}
}
