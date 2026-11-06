package main

import "testing"

func TestSingleQuotes(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"echo 'hello world'",
			"echo 'hello     world'",
			"echo 'hello\\nworld'",
			"echo 'hello \"world\"'",
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
			"$ hello \"world\"",
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})
}
