package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type cd struct {
	shell
}

func (c *cd) setStdout(stdout io.Writer) {
	c.stdout = stdout
}

func (c *cd) setStderr(stderr io.Writer) {
	c.stderr = stderr
}

func (c cd) execute(args []string) {
	if len(args) > 1 {
		fmt.Fprintln(c.stderr, "cd: too many arguments")
		return
	}

	path := "~"
	if len(args) != 0 {
		path = args[0]
	}

	if strings.HasPrefix(path, "~") {
		home := getHomeDirectory()
		path = strings.Replace(path, "~", home, 1)
	}

	if !strings.HasPrefix(path, "~") && !strings.HasPrefix(path, ".") && !strings.HasPrefix(path, "/") {
		path = getWorkingDirectory() + "/" + path
	}

	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprint(c.stderr, "cd: ", path, ": No such file or directory\n")
		return
	}

	if !info.IsDir() {
		fmt.Fprint(c.stderr, "cd: ", path, ": Not a directory\n")
		return
	}

	os.Chdir(path)
}
