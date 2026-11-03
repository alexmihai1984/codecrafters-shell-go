package main

import (
	"testing"
)

func TestInvalidCommand(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		"foobar",
		"exit",
	})

	expectedStdout := "$ $ "
	if stdout != expectedStdout {
		t.Errorf("expected stdout %q, got %q", expectedStdout, stdout)
	}

	expectedStderr := "foobar: command not found\n"
	if stderr != expectedStderr {
		t.Errorf("expected stderr %q, got %q", expectedStderr, stderr)
	}
}
