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

	expectedStdout := "$ blueberry grape mango\n$ pear mango\n$ "
	if stdout != expectedStdout {
		t.Errorf(
			"unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout,
			stdout,
		)
	}

	if len(stderr) != 0 {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}
