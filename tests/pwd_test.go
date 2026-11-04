package main

import (
	"testing"
)

func TestPwd(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		"pwd",
		"exit",
	})

	expectedStdout := "$ " + getWd(t) + "\n$ "
	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}
