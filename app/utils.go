package main

import (
	"fmt"
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

func getWorkingDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		os.Exit(1)
	}

	return pwd
}
