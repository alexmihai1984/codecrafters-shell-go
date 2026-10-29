package main

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func findInPath(cmd string) string {
	path := os.Getenv("PATH")
	folders := filepath.SplitList(path)

	for _, folder := range folders {
		filepath := filepath.Join(folder, cmd)

		err := unix.Access(filepath, unix.X_OK)
		if err == nil {
			return filepath
		}
	}

	return ""
}
