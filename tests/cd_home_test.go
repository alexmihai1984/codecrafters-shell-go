package main

import (
	"os"
	"testing"
)

func TestCdHomeDirectory(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	tmpDir := t.TempDir()

	stdout, stderr := runShell(t, []string{
		"cd " + tmpDir,
		"cd",
		"pwd",
		"cd " + tmpDir,
		"cd ~",
		"pwd",
		"exit",
	})

	expectStdout(
		t,
		stdout,
		[]string{
			"$ $ $ " + homeDir,
			"$ $ $ " + homeDir,
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})
}
