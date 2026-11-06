package main

import (
	"testing"
)

func TestREPL(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"invalid_command_1",
			"invalid_command_2",
			"invalid_command_3",
			"invalid_command_4",
			"exit",
		},
	)

	expectStdout(
		t,
		stdout,
		[]string{
			"$ $ $ $ $ ",
		},
	)

	expectStderr(
		t,
		stderr,
		[]string{
			"invalid_command_1: command not found",
			"invalid_command_2: command not found",
			"invalid_command_3: command not found",
			"invalid_command_4: command not found\n",
		},
	)
}
