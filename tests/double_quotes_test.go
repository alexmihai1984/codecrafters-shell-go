package main

import "testing"

func TestDoubleQuotes(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		"echo \"hello world\"",
		"echo \"hello     world\"",
		"echo \"hello\\nworld\"",
		"exit",
	})

	expectedStdout := "$ hello world\n" +
		"$ hello     world\n" +
		"$ hello\\nworld\n" +
		"$ "

	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}
