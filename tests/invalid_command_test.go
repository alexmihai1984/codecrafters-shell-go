package main

import (
	"testing"
)

func TestInvalidCommand(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"foobar",
			"exit",
		},
	)

	expectStdout(
		t,
		stdout,
		[]string{
			"$ $ ",
		},
	)

	expectStderr(
		t,
		stderr,
		[]string{
			"foobar: command not found\n",
		},
	)
}
