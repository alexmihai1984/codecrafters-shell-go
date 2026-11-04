package main

import "testing"

func TestSingleQuotes(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		"echo 'hello world'",
		"echo 'hello     world'",
		"echo 'hello\\nworld'",
		"echo 'hello \"world\"'",
		"exit",
	})

	expectedStdout := "$ hello world\n" +
		"$ hello     world\n" +
		"$ hello\\nworld\n" +
		"$ hello \"world\"\n" +
		"$ "

	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}
