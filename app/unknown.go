package main

import (
	"fmt"
)

type unknown struct {
	verb string
}

func (u unknown) execute(args []string) {
	fmt.Print(u.verb, ": command not found\n")
}
