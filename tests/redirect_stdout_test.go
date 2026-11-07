package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStdoutRedirection(t *testing.T) {
	tests := []struct {
		name        string
		redirection string
	}{
		{
			name:        "greater than with spaces",
			redirection: " > ",
		},
		{
			name:        "greater than without spaces",
			redirection: ">",
		},
		{
			name:        "1 greater than with spaces",
			redirection: " 1> ",
		},
		{
			name:        "1 greater than without space before filename",
			redirection: " 1>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			outputFile := filepath.Join(dir, "output.txt")

			command := "/bin/echo hello" + tt.redirection + outputFile

			stdout, stderr := runShell(t, []string{
				command,
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

			expectFileContent(t, outputFile, []string{"hello\n"})
		})
	}
}

func TestStdoutRedirectionWithTrailingOneInput(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")

	stdout, stderr := runShell(t, []string{
		"echo hello1>" + outputFile,
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

	expectFileContent(t, outputFile, []string{"hello1\n"})
}

func TestStdoutRedirectionWhenCommandFails(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")
	missingFile := filepath.Join(dir, "does-not-exist")

	if err := os.WriteFile(outputFile, []byte("old contents"), 0644); err != nil {
		t.Fatalf("failed to create output file: %v", err)
	}

	stdout, stderr := runShell(t, []string{
		"/bin/cat " + missingFile + " > " + outputFile,
		"exit",
	})

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
			"cat: " + missingFile + ": No such file or directory\n",
		},
	)

	expectFileContent(t, outputFile, []string{})
}

func TestStderrRedirectImmediatelyAfterStdoutRedirect(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")

	stdout, stderr := runShell(t, []string{
		"echo hello > 2> " + outputFile,
		"exit",
	})

	expectStdout(
		t,
		stdout,
		[]string{
			"$ nsh: syntax error near unexpected token `2'",
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})

	expectNoFile(t, outputFile)
}
