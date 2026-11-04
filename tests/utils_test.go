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

func expectStdout(t *testing.T, stdout string, lines []string) {
	expectedStdout := strings.Join(lines, "\n")

	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q", expectedStdout, stdout)
	}
}

func expectStderr(t *testing.T, stderr string, lines []string) {
	expectedStderr := strings.Join(lines, "\n")

	if stderr != expectedStderr {
		t.Errorf("unexpected stderr:\nexpected: %q\ngot:      %q", expectedStderr, stderr)
	}
}

func getWd(t *testing.T) string {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	return dir
}
