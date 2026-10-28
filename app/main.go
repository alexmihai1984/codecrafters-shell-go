package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Print("$ ")
	cmd, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return;
	}

	cmd = strings.TrimSpace(cmd)

	fmt.Println(cmd, ": command not found")
}
