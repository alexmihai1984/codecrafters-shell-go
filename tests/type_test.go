package main

import (
	"testing"
)

func TestTypeBuiltin(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"type echo",
			"type exit",
			"type type",
			"exit",
		},
	)

	expectStdout(
		t,
		stdout,
		[]string{
			"$ echo is a shell builtin",
			"$ exit is a shell builtin",
			"$ type is a shell builtin",
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})
}

func TestTypeNonexistent(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"type nonexistent_command",
			"exit",
		},
	)

	expectStdout(
		t,
		stdout,
		[]string{
			"$ nonexistent_command: not found\n$ ",
		},
	)

	expectStderr(t, stderr, []string{})
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

	expectStdout(
		t,
		stdout,
		[]string{
			"$ ls is /usr/bin/ls",
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})
}
