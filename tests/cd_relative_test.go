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

	expectStdout(
		t,
		stdout,
		[]string{
			"$ $ $ " + firstDir + "\n" +
				"$ $ " + secondDir + "\n" +
				"$ $ " + firstDir + "\n" +
				"$ ",
		},
	)

	expectStderr(t, stderr, []string{})
}

func TestCdRelativeNonexistentPath(t *testing.T) {
	invalidDir := "./this/directory/does/not/exist"

	stdout, stderr := runShell(t, []string{
		"cd " + invalidDir,
		"pwd",
		"exit",
	})

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
