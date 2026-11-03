package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestInvalidCommand(t *testing.T) {
	cmd := exec.Command("../your_program.sh")
	commands := []string{
		"foobar",
		"exit",
	}

	cmd.Stdin = strings.NewReader(strings.Join(commands, "\n") + "\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("program failed: %v", err)
	}

	expectedStdout := "$ $ "
	if stdout.String() != expectedStdout {
		t.Errorf("expected stdout %q, got %q", expectedStdout, stdout.String())
	}

	expectedStderr := "foobar: command not found\n"
	if stderr.String() != expectedStderr {
		t.Errorf("expected stderr %q, got %q", expectedStderr, stderr.String())
	}
}
