package main

import "testing"

func TestCdAbsolutePath(t *testing.T) {
	firstDir := t.TempDir()
	secondDir := t.TempDir()

	stdout, stderr := runShell(
		t,
		[]string{
			"cd " + firstDir,
			"pwd",
			"cd " + secondDir,
			"pwd",
			"exit",
		},
	)

	expectStdout(
		t,
		stdout,
		[]string{
			"$ $ " + firstDir,
			"$ $ " + secondDir,
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})
}

func TestCdAbsoluteNonexistentPath(t *testing.T) {
	invalidDir := "/this/directory/does/not/exist"

	stdout, stderr := runShell(
		t,
		[]string{
			"cd " + invalidDir,
			"pwd",
			"exit",
		},
	)

	expectStdout(
		t,
		stdout,
		[]string{
			"$ $ " + getWd(t),
			"$ ",
		},
	)

	expectStderr(
		t,
		stderr,
		[]string{
			"cd: " + invalidDir + ": No such file or directory\n",
		},
	)
}
