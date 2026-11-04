package main

import "testing"

func TestCdAbsolutePath(t *testing.T) {
	firstDir := t.TempDir()
	secondDir := t.TempDir()

	stdout, stderr := runShell(t, []string{
		"cd " + firstDir,
		"pwd",
		"cd " + secondDir,
		"pwd",
		"exit",
	})

	expectedStdout := "$ $ " + firstDir + "\n$ $ " + secondDir + "\n$ "

	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}

func TestCdAbsoluteNonexistentPath(t *testing.T) {
	invalidDir := "/this/directory/does/not/exist"

	stdout, stderr := runShell(t, []string{
		"cd " + invalidDir,
		"pwd",
		"exit",
	})

	expectedStdout := "$ $ " + getWd(t) + "\n$ "
	expectedStderr := "cd: " + invalidDir + ": No such file or directory\n"

	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout, stdout)
	}

	if stderr != expectedStderr {
		t.Errorf("unexpected stderr:\nexpected: %q\ngot:      %q",
			expectedStderr, stderr)
	}
}
