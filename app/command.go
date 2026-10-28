package main

import "strings"

type command struct {
	verb string
	args []string
}

func parseCmd(cmd string) command {
	parts := strings.Fields(cmd)

	if len(parts) == 0 {
		return command{}
	}

	return command{
		parts[0],
		parts[1:],
	}
}
