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
		if !handleCmd(cmd) {
			return
		}
	}
}

func handleCmd(cmd string) bool {
	var command command = parseCmd(cmd)
	if command.verb == "exit" {
		return false
	}

	if command.verb == "echo" {
		fmt.Println(strings.Join(command.args, " "))
	} else {
		fmt.Print(cmd, ": command not found\n")
	}

	return true
}
