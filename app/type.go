package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
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
		} else if !findInPath(cmd) {
			fmt.Print(cmd, ": not found\n")
		}
	}
}

func findInPath(cmd string) bool {
	path := os.Getenv("PATH")
	folders := filepath.SplitList(path)

	for _, folder := range folders {
		filepath := filepath.Join(folder, cmd)

		err := unix.Access(filepath, unix.X_OK)
		if err == nil {
			fmt.Println(cmd, "is", filepath)
			return true
		}
	}

	return false
}
