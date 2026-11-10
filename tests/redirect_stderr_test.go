package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStderrRedirection(t *testing.T) {
	tests := []struct {
		name        string
		redirection string
	}{
		{
			name:        "2 greater than with spaces",
			redirection: " 2> ",
		},
		{
			name:        "2 greater than without space before filename",
			redirection: " 2>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			outputFile := filepath.Join(dir, "output.txt")

			stdout, stderr := runShell(t, []string{
				"cat does-not-exist" + tt.redirection + outputFile,
				"exit",
			})

			expectStdout(
				t,
				stdout,
				[]string{
					"$ $ ",
				},
			)

			expectStderr(t, stderr, []string{})

			expectFileContent(
				t,
				outputFile,
				[]string{
					"cat: does-not-exist: No such file or directory\n",
				},
			)
		})
	}
}

func TestStderrRedirectionWithTrailingTwoInput(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")

	stdout, stderr := runShell(t, []string{
		"echo hello2>" + outputFile,
		"exit",
	})

	expectStdout(
		t,
		stdout,
		[]string{
			"$ $ ",
		},
	)

	expectStderr(t, stderr, []string{})

	expectFileContent(t, outputFile, []string{"hello2\n"})
}

func TestStderrRedirectionWhenCommandSucceeds(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")

	if err := os.WriteFile(outputFile, []byte("old contents"), 0644); err != nil {
		t.Fatalf("failed to create output file: %v", err)
	}

	stdout, stderr := runShell(t, []string{
		"echo hello 2> " + outputFile,
		"exit",
	})

	expectStdout(
		t,
		stdout,
		[]string{
			"$ hello",
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})

	expectFileContent(t, outputFile, []string{})
}

func TestStdoutRedirectImmediatelyAfterStderrRedirect(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")

	stdout, stderr := runShell(t, []string{
		"echo hello 2> > " + outputFile,
		"exit",
	})

	expectStdout(
		t,
		stdout,
		[]string{
			"$ nsh: syntax error near unexpected token `>'",
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})

	expectNoFile(t, outputFile)
}

func TestMultipleStderrRedirections(t *testing.T) {
	dir := t.TempDir()
	outputFile1 := filepath.Join(dir, "output1.txt")
	outputFile2 := filepath.Join(dir, "output2.txt")

	stdout, stderr := runShell(
		t,
		[]string{
			"cat does-not-exits 2> " + outputFile1 + " 2> " + outputFile2,
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

	expectStderr(t, stderr, []string{})

	expectFileContent(t, outputFile1, []string{})
	expectFileContent(t, outputFile2, []string{"cat: does-not-exits: No such file or directory\n"})
}

func TestRedirectStderrWithoutCommand(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")

	stdout, stderr := runShell(
		t,
		[]string{
			"2> " + outputFile,
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

	expectStderr(t, stderr, []string{})

	expectFileContent(t, outputFile, []string{})
}
