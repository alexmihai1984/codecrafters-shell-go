package main

import (
	"strings"
	"testing"
)

func TestREPL(t *testing.T) {
	commands := []string{
		"invalid_command_1",
		"invalid_command_2",
		"invalid_command_3",
		"invalid_command_4",
		"exit",
	}
	stdout, stderr := runShell(t, commands)

	expectedStdout := strings.Repeat("$ ", 5)

	var expectedStderr strings.Builder
	for _, command := range commands[:4] {
		expectedStderr.WriteString(command + ": command not found\n")
	}

	if stdout != expectedStdout {
		t.Errorf("expected stdout %q, got %q", expectedStdout, stdout)
	}

	if stderr != expectedStderr.String() {
		t.Errorf("expected stderr %q, got %q", expectedStderr.String(), stderr)
	}
}
