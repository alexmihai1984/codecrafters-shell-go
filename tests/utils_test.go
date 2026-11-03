package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func runShell(t *testing.T, commands []string) (string, string) {
	return runShellWithPath(t, commands, "")
}

func runShellWithPath(t *testing.T, commands []string, path string) (string, string) {
	cmd := exec.Command("../your_program.sh")
	if path != "" {
		env := os.Environ()

		for i, entry := range env {
			if strings.HasPrefix(entry, "PATH=") {
				env[i] = "PATH=" + path
				break
			}
		}

		cmd.Env = env
	}

	cmd.Stdin = strings.NewReader(strings.Join(commands, "\n") + "\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("shell failed: %v\nstdout: %q\nstderr: %q", err, stdout.String(), stderr.String())
	}

	return stdout.String(), stderr.String()
}
