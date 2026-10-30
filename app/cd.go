package main

import (
	"fmt"
	"os"
	"strings"
)

type cd struct{}

func (c cd) execute(args []string) {
	if len(args) > 1 {
		fmt.Println("cd: too many arguments")
		return
	}

	path := args[0]

	if strings.HasPrefix(path, "~") {
		home := getHomeDirectory()
		path = strings.Replace(path, "~", home, 1)
	}

	if !strings.HasPrefix(path, "~") && !strings.HasPrefix(path, ".") && !strings.HasPrefix(path, "/") {
		path = getWorkingDirectory() + "/" + path
	}

	info, err := os.Stat(path)
	if err != nil {
		fmt.Print("cd: ", path, ": No such file or directory\n")
		return
	}

	if !info.IsDir() {
		fmt.Print("cd: ", path, ": Not a directory\n")
		return
	}

	os.Chdir(path)
}
