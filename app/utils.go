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

func findAbsolutePath(cmd string) string {
	err := unix.Access(cmd, unix.X_OK)
	if err != nil {
		return ""
	}

	return cmd
}

func getWorkingDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		os.Exit(1)
	}

	return pwd
}

func getHomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("cd: Could not get home directory")
		os.Exit(1)
	}

	return home
}
