package main

import (
	"testing"
)

func TestEchoBuiltin(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		"echo blueberry grape mango",
		"echo pear mango",
		"exit",
	})

	expectStdout(
		t,
		stdout,
		[]string{
			"$ blueberry grape mango",
			"$ pear mango\n$ ",
		},
	)

	expectStderr(t, stderr, []string{})
}
