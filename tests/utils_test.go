package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func runShell(t *testing.T, commands []string) (string, string) {
	cmd := exec.Command("../your_program.sh")

	cmd.Stdin = strings.NewReader(strings.Join(commands, "\n") + "\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("shell failed: %v", err)
	}

	return stdout.String(), stderr.String()
}
