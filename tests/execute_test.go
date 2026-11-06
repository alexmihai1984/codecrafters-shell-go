package main

import (
	"testing"
)

func TestExecuteProgram(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"ls",
			"exit",
		},
	)

	expectStdoutContains(
		t,
		stdout,
		"execute_test.go",
	)

	expectStderr(t, stderr, []string{})
}

func TestExecuteProgramWithArguments(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"printf 'hello from external program\\n'",
			"exit",
		},
	)

	expectStdout(
		t,
		stdout,
		[]string{
			"$ hello from external program",
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})
}

func TestExecuteNonexistentProgram(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"nonexistent_command",
			"exit",
		},
	)

	expectStdout(
		t,
		stdout,
		[]string{
			"$ $ ",
		},
	)

	expectStderr(
		t,
		stderr,
		[]string{
			"nonexistent_command: command not found\n",
		},
	)
}

func TestAbsoluteExecutablePath(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		`/bin/echo hello`,
		"exit",
	})

	expectStdout(t, stdout, []string{
		"$ hello",
		"$ ",
	})

	expectStderr(
		t,
		stderr,
		[]string{},
	)
}
