package main

import "io"

type executable interface {
	execute(args []string)
	setStdout(stdout io.Writer)
	setStderr(stderr io.Writer)
}
