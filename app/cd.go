package main

import (
	"fmt"
	"os"
)

type cd struct{}

func (c cd) execute(args []string) {
	if len(args) > 1 {
		fmt.Println("cd: too many arguments")
		return
	}

	info, err := os.Stat(args[0])
	if err != nil {
		fmt.Print("cd: ", args[0], ": No such file or directory\n")
		return
	}

	if !info.IsDir() {
		fmt.Print("cd: ", args[0], ": Not a directory\n")
		return
	}

	os.Chdir(args[0])
}
