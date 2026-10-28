package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		cmd = strings.TrimSpace(cmd)
		if !processCmd(cmd) {
			return
		}

		fmt.Print(cmd, ": command not found\n")
	}
}

func processCmd(cmd string) bool {
	if cmd == "exit" {
		return false
	}

	return true
}
