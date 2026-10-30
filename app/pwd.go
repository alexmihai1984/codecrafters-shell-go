package main

import (
	"fmt"
)

type pwd struct{}

func (p pwd) execute(args []string) {
	fmt.Println(getWorkingDirectory())
}
