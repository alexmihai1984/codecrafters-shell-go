package main

import (
	"strings"
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

	if !strings.Contains(stdout, "execute_test.go") {
		t.Errorf("expected ls output to contain execute_test.go, got: %q", stdout)
	}

	if len(stderr) != 0 {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}

func TestExecuteProgramWithArguments(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"printf 'hello from external program\\n'",
			"exit",
		},
	)

	expectedStdout := "$ hello from external program\n$ "
	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q", expectedStdout, stdout)
	}

	if len(stderr) != 0 {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}

func TestExecuteNonexistentProgram(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"nonexistent_command",
			"exit",
		},
	)

	expectedStdout := "$ $ "
	expectedStderr := "nonexistent_command: command not found\n"

	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout, stdout)
	}

	if stderr != expectedStderr {
		t.Errorf("unexpected stderr:\nexpected: %q\ngot:      %q",
			expectedStderr, stderr)
	}
}
