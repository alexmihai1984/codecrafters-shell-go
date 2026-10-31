package main

import (
	"fmt"
	"io"
)

type pwd struct {
	shell
}

func (p *pwd) setStdout(stdout io.Writer) {
	p.stdout = stdout
}

func (p *pwd) setStderr(stderr io.Writer) {
	p.stderr = stderr
}

func (p pwd) execute(args []string) {
	fmt.Println(getWorkingDirectory())
}
