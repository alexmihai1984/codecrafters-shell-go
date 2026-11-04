package main

import (
	"os"
	"path/filepath"
	"testing"
)

// todo refactor these tests
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

			if stdout != "$ $ " {
				t.Errorf("unexpected stdout: %q", stdout)
			}

			if stderr != "" {
				t.Errorf("unexpected stderr: %q", stderr)
			}

			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("failed to read redirected file: %v", err)
			}

			if string(content) != "hello\n" {
				t.Errorf("expected file content %q, got %q", "hello\n", string(content))
			}
		})
	}
}

func TestStdoutRedirectionWithTrailingOne(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")

	stdout, stderr := runShell(t, []string{
		"echo hello1>" + outputFile,
		"exit",
	})

	if stdout != "$ $ " {
		t.Errorf("unexpected stdout: %q", stdout)
	}

	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read redirected file: %v", err)
	}

	if string(content) != "hello1\n" {
		t.Errorf("expected %q, got %q", "hello1\n", string(content))
	}
}

func TestStdoutRedirectionWhenCommandFails(t *testing.T) {
	dir := t.TempDir()
	outputFile := filepath.Join(dir, "output.txt")
	missingFile := filepath.Join(dir, "does-not-exist")

	// Prepopulate the output file to verify that redirection truncates it.
	if err := os.WriteFile(outputFile, []byte("old contents"), 0644); err != nil {
		t.Fatalf("failed to create output file: %v", err)
	}

	stdout, stderr := runShell(t, []string{
		"/bin/cat " + missingFile + " > " + outputFile,
		"exit",
	})

	if stdout != "$ $ " {
		t.Errorf("unexpected stdout: %q", stdout)
	}

	expectedStderr := "cat: " + missingFile + ": No such file or directory\n"
	if stderr != expectedStderr {
		t.Errorf("unexpected stderr:\nexpected: %q\ngot:      %q", expectedStderr, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read redirected file: %v", err)
	}

	if len(content) != 0 {
		t.Errorf("expected empty file, got %q", string(content))
	}
}
