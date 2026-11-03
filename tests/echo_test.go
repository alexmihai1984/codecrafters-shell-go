package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestEchoBuiltin(t *testing.T) {
	cmd := exec.Command("../your_program.sh")
	commands := []string{
		"echo blueberry grape mango",
		"echo pear mango",
		"exit",
	}

	cmd.Stdin = strings.NewReader(strings.Join(commands, "\n") + "\n")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("shell failed: %v", err)
	}

	expectedStdout := "$ blueberry grape mango\n$ pear mango\n$ "
	if stdout.String() != expectedStdout {
		t.Errorf(
			"unexpected stdout:\nexpected: %q\ngot:      %q",
			expectedStdout,
			stdout.String(),
		)
	}

	if stderr.Len() != 0 {
		t.Errorf("unexpected stderr: %q", stderr.String())
	}
}
