package main

import (
	"path/filepath"
	"testing"
)

func TestQuotedExecutableName(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name        string
		exeName     string
		contentFile string
		shellInput  string
	}{
		{
			name:        "exe with space",
			exeName:     `exe  with  space`,
			contentFile: "content_f1.txt",
			shellInput:  "'" + filepath.Join(dir, `exe  with  space`) + "' " + filepath.Join(dir, "content_f1.txt"),
		},
		{
			name:        "exe with double quotes",
			exeName:     `exe with "quotes"`,
			contentFile: "content_f2.txt",
			shellInput:  "'" + filepath.Join(dir, `exe with "quotes"`) + "' " + filepath.Join(dir, "content_f2.txt"),
		},
		{
			name:        "exe with single quotes",
			exeName:     `exe with 'single quotes'`,
			contentFile: "content_f3.txt",
			shellInput:  "\"" + filepath.Join(dir, `exe with 'single quotes'`) + "\" " + filepath.Join(dir, "content_f3.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exeFilePath := filepath.Join(dir, tt.exeName)
			contentFilePath := filepath.Join(dir, tt.contentFile)

			createFile(
				t,
				exeFilePath,
				[]string{"#!/bin/sh\nexec /bin/cat \"$@\"\n"},
				0755,
			)

			createFile(
				t,
				contentFilePath,
				[]string{"some content\n"},
				0644,
			)

			stdout, stderr := runShell(t, []string{
				tt.shellInput,
				"exit",
			})

			expectStdout(
				t,
				stdout,
				[]string{"$ some content", "$ "},
			)

			expectStderr(t, stderr, []string{})
		})
	}
}
