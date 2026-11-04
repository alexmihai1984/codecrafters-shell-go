package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCdRelativePath(t *testing.T) {
	rootDir := t.TempDir()
	firstDir := filepath.Join(rootDir, "first")
	secondDir := filepath.Join(firstDir, "second")

	if err := os.MkdirAll(secondDir, 0755); err != nil {
		t.Fatalf("failed to create directories: %v", err)
	}

	stdout, stderr := runShell(t, []string{
		"cd " + rootDir,
		"cd ./first",
		"pwd",
		"cd second",
		"pwd",
		"cd ..",
		"pwd",
		"exit",
	})

	expectedStdout := "$ $ $ " + firstDir + "\n" +
		"$ $ " + secondDir + "\n" +
		"$ $ " + firstDir + "\n" +
		"$ "

	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout, stdout)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}

func TestCdRelativeNonexistentPath(t *testing.T) {
	invalidDir := "./this/directory/does/not/exist"

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
