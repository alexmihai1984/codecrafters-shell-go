package main

import "io"

type noop struct {
	shell
}

func (n *noop) setStdout(stdout io.Writer) {
	n.stdout = stdout
}

func (n *noop) setStderr(stderr io.Writer) {
	n.stderr = stderr
}

func (n noop) execute(args []string) {

}
