package main

import (
	"path/filepath"
	"testing"
)

func TestStdoutAppendRedirection(t *testing.T) {
	tests := []struct {
		name        string
		redirection string
	}{
		{
			name:        "greater than with spaces",
			redirection: " >> ",
		},
		{
			name:        "greater than without spaces",
			redirection: ">>",
		},
		{
			name:        "1 greater than with spaces",
			redirection: " 1>> ",
		},
		{
			name:        "1 greater than without space before filename",
			redirection: " 1>>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			outputFile := filepath.Join(dir, "output.txt")

			createFile(
				t,
				outputFile,
				[]string{"some content\n"},
				0644,
			)

			command := "echo hello" + tt.redirection + outputFile

			stdout, stderr := runShell(
				t,
				[]string{
					command,
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

			expectFileContent(
				t,
				outputFile,
				[]string{
					"some content",
					"hello\n",
				},
			)
		})
	}
}

func TestStdoutAppendRedirectionWithTrailingOneInput(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")

	createFile(
		t,
		outputFile,
		[]string{"some content\n"},
		0644,
	)

	stdout, stderr := runShell(
		t,
		[]string{
			"echo hello1>>" + outputFile,
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

	expectFileContent(
		t,
		outputFile,
		[]string{
			"some content",
			"hello1\n",
		},
	)
}

func TestStdoutAppendRedirectionWhenCommandFails(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")
	missingFile := filepath.Join(dir, "does-not-exist")

	createFile(
		t,
		outputFile,
		[]string{"old contents\n"},
		0644,
	)

	stdout, stderr := runShell(t, []string{
		"/bin/cat " + missingFile + " >> " + outputFile,
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

	expectFileContent(t, outputFile, []string{"old contents\n"})
}

func TestStderrAppendRedirectImmediatelyAfterStdoutRedirect(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")

	stdout, stderr := runShell(
		t,
		[]string{
			"echo hello >> 2>> " + outputFile,
			"exit",
		},
	)

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

func TestMultipleStdoutAppendRedirects(t *testing.T) {
	dir := t.TempDir()
	outputFile1 := filepath.Join(dir, "output1.txt")
	outputFile2 := filepath.Join(dir, "output2.txt")

	stdout, stderr := runShell(
		t,
		[]string{
			"echo hello >> " + outputFile1 + " >> " + outputFile2,
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
	expectFileContent(t, outputFile2, []string{"hello\n"})
}

func TestStdoutAppendRedirectionWithoutCommand(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")

	stdout, stderr := runShell(
		t,
		[]string{
			">> " + outputFile,
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
