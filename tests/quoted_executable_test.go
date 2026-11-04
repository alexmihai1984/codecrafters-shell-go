package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestQuotedExecutableNames(t *testing.T) {
	binDir := t.TempDir()

	executableNames := []string{
		`exe  with  space`,
		`exe with "quotes"`,
		`exe with 'single quotes'`,
	}

	for _, name := range executableNames {
		path := filepath.Join(binDir, name)
		script := "#!/bin/sh\nexec /bin/cat \"$@\"\n"

		if err := os.WriteFile(path, []byte(script), 0755); err != nil {
			t.Fatalf("failed to create executable: %v", err)
		}
	}

	files := []string{"f1", "f2", "f3"}
	contents := []string{
		"mango raspberry.",
		"orange raspberry.",
		"pineapple orange.",
	}

	for i, name := range files {
		if err := os.WriteFile(filepath.Join(binDir, name), []byte(contents[i]), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	stdout, stderr := runShell(t, []string{
		"'" + filepath.Join(binDir, `exe  with  space`) + "' " + filepath.Join(binDir, "f1"),
		"'" + filepath.Join(binDir, `exe with "quotes"`) + "' " + filepath.Join(binDir, "f2"),
		"\"" + filepath.Join(binDir, `exe with 'single quotes'`) + "\" " + filepath.Join(binDir, "f3"),
		"exit",
	})

	expectedStdout := "$ mango raspberry." +
		"$ orange raspberry." +
		"$ pineapple orange." +
		"$ "

	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q", expectedStdout, stdout)
	}
	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}
