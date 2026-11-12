package main

import (
	"errors"
	"fmt"
)

type interpreter struct {
	parts []string
}

func buildInterpreter(parts []string) interpreter {
	return interpreter{
		parts: parts,
	}
}

func (interpreter *interpreter) interpret() command {
	if len(interpreter.parts) == 0 {
		return command{
			verb:       "noop",
			args:       nil,
			redirects:  nil,
			executable: &noop{},
		}
	}

	res := command{}
	previousOneWasRedirect := false
	for i, token := range interpreter.parts {
		if isRedirect(token) {
			previousOneWasRedirect = true
			redirectType := interpreter.getRedirectType(token)
			file, err := interpreter.getRedirectFile(i + 1)
			if err != nil {
				fmt.Println(err.Error())
				return command{
					verb:       "noop",
					args:       nil,
					redirects:  nil,
					executable: &noop{},
				}
			}

			res.redirects = append(res.redirects, redirect{redirectType: redirectType, file: file})
		} else if res.verb == "" {
			if previousOneWasRedirect {
				previousOneWasRedirect = false
			} else {
				res.verb = token
			}
		} else {
			if previousOneWasRedirect {
				previousOneWasRedirect = false
			} else {
				res.args = append(res.args, token)
			}
		}
	}

	executable, ok := builtins[res.verb]
	if !ok {
		res.executable = &program{
			name: res.verb,
		}

		return res
	}

	res.executable = executable

	return res
}

func (interpreter *interpreter) getRedirectFile(i int) (string, error) {
	if len(interpreter.parts) <= i {
		return "", errors.New("nsh: syntax error near unexpected token `newline'")
	}

	if isRedirect(interpreter.parts[i]) {
		return "", errors.New("nsh: syntax error near unexpected token `" + string(interpreter.parts[i][0]) + "'")
	}

	return interpreter.parts[i], nil
}

func (interpreter *interpreter) getRedirectType(s string) redirectType {
	if isStdoutRedirect(s) {
		return Stdout
	} else if isStderrRedirect(s) {
		return Stderr
	} else {
		return None
	}
}
func isRedirect(s string) bool {
	return isStdoutRedirect(s) || isStderrRedirect(s)
}

func isStdoutRedirect(s string) bool {
	return s == ">" || s == "1>"
}

func isStderrRedirect(s string) bool {
	return s == "2>"
}
