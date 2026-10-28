package main

import (
	"fmt"
	"strings"
)

type echo struct{}

func (e echo) execute(args []string) {
	fmt.Println(strings.Join(args, " "))
}
