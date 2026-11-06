package main

import "testing"

func TestDoubleQuotes(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"echo \"hello world\"",
			"echo \"hello     world\"",
			"echo \"hello\\nworld\"",
			"exit",
		},
	)

	expectStdout(
		t,
		stdout,
		[]string{
			"$ hello world",
			"$ hello     world",
			"$ hello\\nworld",
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})
}
