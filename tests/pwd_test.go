package main

import (
	"os"
	"testing"
)

func TestPwd(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		"pwd",
		"exit",
	})

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	expectedStdout := "$ " + dir + "\n$ "
	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}
