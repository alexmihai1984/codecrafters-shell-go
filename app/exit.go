package main

import (
	"io"
	"os"
)

type exit struct {
	shell
}

func (e *exit) setStdout(stdout io.Writer) {
	e.stdout = stdout
}

func (e *exit) setStderr(stderr io.Writer) {
	e.stderr = stderr
}

func (e exit) execute(args []string) {
	os.Exit(0)
}
