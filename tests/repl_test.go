package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestREPL(t *testing.T) {
	cmd := exec.Command("../your_program.sh")
	commands := []string{
		"invalid_command_1",
		"invalid_command_2",
		"invalid_command_3",
		"invalid_command_4",
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

	expectedStdout := strings.Repeat("$ ", 5)

	var expectedStderr strings.Builder
	for _, command := range commands[:4] {
		expectedStderr.WriteString(command + ": command not found\n")
	}

	if stdout.String() != expectedStdout {
		t.Errorf("expected stdout %q, got %q", expectedStdout, stdout.String())
	}

	if stderr.String() != expectedStderr.String() {
		t.Errorf("expected stderr %q, got %q", expectedStderr.String(), stderr.String())
	}
}
