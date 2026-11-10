package main

import (
	"io"
)

type shell struct {
	stdout io.Writer
	stderr io.Writer
}
