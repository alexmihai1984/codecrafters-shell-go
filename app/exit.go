package main

import (
	"os"
)

type exit struct{}

func (e exit) execute(args []string) {
	os.Exit(0)
}
