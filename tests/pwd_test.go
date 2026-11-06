package main

import (
	"testing"
)

func TestPwd(t *testing.T) {
	stdout, stderr := runShell(
		t,
		[]string{
			"pwd",
			"exit",
		},
	)

	expectStdout(
		t,
		stdout,
		[]string{
			"$ " + getWd(t),
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})
}
