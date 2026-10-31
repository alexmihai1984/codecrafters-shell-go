package main

import (
	"fmt"
	"io"
	"strings"
)

type echo struct {
	shell
}

func (e *echo) setStdout(stdout io.Writer) {
	e.stdout = stdout
}

func (e *echo) setStderr(stderr io.Writer) {
	e.stderr = stderr
}

func (e echo) execute(args []string) {
	fmt.Fprintln(e.stdout, strings.Join(args, " "))
}
