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
		handleCmd(cmd)
	}
}

func handleCmd(cmd string) {
	var command command = parseCmd(cmd)

	command.execute()
}
