package main

import (
	"testing"
)

func TestTypeBuiltin(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		"type echo",
		"type exit",
		"type type",
		"exit",
	})

	expectedStdout := "$ echo is a shell builtin\n" +
		"$ exit is a shell builtin\n" +
		"$ type is a shell builtin\n" +
		"$ "

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

func TestTypeNonexistent(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		"type nonexistent_command",
		"exit",
	})

	expectedStdout := "$ nonexistent_command: not found\n$ "
	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q", expectedStdout, stdout)
	}

	if len(stderr) != 0 {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}

func TestTypeProgram(t *testing.T) {
	stdout, stderr := runShellWithPath(
		t,
		[]string{
			"type ls",
			"exit",
		},
		"/usr/bin:/usr/local/bin:/bin",
	)

	expectedStdout := "$ ls is /usr/bin/ls\n$ "
	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q", expectedStdout, stdout)
	}

	if len(stderr) != 0 {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}
