package main

import (
	"fmt"
	"io"
)

type unknown struct {
	shell
	verb string
}

func (u *unknown) setStdout(stdout io.Writer) {
	u.stdout = stdout
}

func (u *unknown) setStderr(stderr io.Writer) {
	u.stderr = stderr
}

func (u unknown) execute(args []string) {
	fmt.Fprint(u.stderr, u.verb, ": command not found\n")
}
