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

	expectedStdout := "$ $ $ " + homeDir + "\n" +
		"$ $ $ " + homeDir + "\n" +
		"$ "

	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}
