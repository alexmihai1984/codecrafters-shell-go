package main

type executable interface {
	execute(args []string)
}
