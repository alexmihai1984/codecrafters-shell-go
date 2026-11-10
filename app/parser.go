package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type parser struct {
	cmd   string
	state parserState
}

type parserState struct {
	parts       []string
	currentWord strings.Builder
}

type parseResult struct {
	ok        bool
	verb      string
	params    []string
	redirects []redirect
}

type redirect struct {
	redirectType redirectType
	file         string
}

type redirectType int

const (
	Stdout redirectType = iota
	Stderr
	StdoutAppend
	StderrAppend
	None
)

func buildParser(cmd string) parser {
	return parser{
		cmd: cmd,
	}
}

func (p *parser) parse() parseResult {
	token := ""
	for i := 0; i < len(p.cmd); {
		token, i = p.readAt(i)
		if token != "" {
			p.state.parts = append(p.state.parts, token)
		}
	}

	if len(p.state.parts) == 0 {
		return parseResult{
			ok: true,
		}
	}

	res := parseResult{
		ok: true,
	}
	previousOneWasRedirect := false
	for i, token := range p.state.parts {
		if p.isRedirect(token) {
			previousOneWasRedirect = true
			redirectType := p.getRedirectType(token)
			file, err := p.getRedirectFile(i + 1)
			if err != nil {
				fmt.Println(err.Error())
				return parseResult{
					ok: false,
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
				res.params = append(res.params, token)
			}
		}
	}

	return res
}

func (p *parser) getRedirectFile(i int) (string, error) {
	if len(p.state.parts) <= i {
		return "", errors.New("nsh: syntax error near unexpected token `newline'")
	}

	if p.isRedirect(p.state.parts[i]) {
		return "", errors.New("nsh: syntax error near unexpected token `" + string(p.state.parts[i][0]) + "'")
	}

	return p.state.parts[i], nil
}

// todo: split interpreter from parser
// todo: parseResult and command are very similar, maybe reconciliate them
func (p *parser) getRedirectType(s string) redirectType {
	if p.isStdoutRedirect(s) {
		return Stdout
	} else if p.isStderrRedirect(s) {
		return Stderr
	} else {
		return None
	}
}

func (p *parser) getErrRedirect() (int, string, error) {
	errRedirectIdx := p.findErrRedirect()
	if errRedirectIdx < 0 {
		return -1, "", nil
	}

	if len(p.state.parts) <= errRedirectIdx+1 {
		return -1, "", errors.New("nsh: syntax error near unexpected token `newline'")
	}

	if p.isRedirect(p.state.parts[errRedirectIdx+1]) {
		return -1, "", errors.New("nsh: syntax error near unexpected token `" + string(p.state.parts[errRedirectIdx+1][0]) + "'")
	}

	return errRedirectIdx, p.state.parts[errRedirectIdx+1], nil
}

func (p *parser) findErrRedirect() int {
	for i, s := range p.state.parts {
		if s == "2>" {
			return i
		}
	}

	return -1
}

func (p *parser) getStdRedirect() (int, string, error) {
	stdRedirectIdx := p.findStdRedirect()
	if stdRedirectIdx < 0 {
		return -1, "", nil
	}

	if len(p.state.parts) <= stdRedirectIdx+1 {
		return -1, "", errors.New("nsh: syntax error near unexpected token `newline'")
	}

	if p.isRedirect(p.state.parts[stdRedirectIdx+1]) {
		return -1, "", errors.New("nsh: syntax error near unexpected token `" + string(p.state.parts[stdRedirectIdx+1][0]) + "'")
	}

	return stdRedirectIdx, p.state.parts[stdRedirectIdx+1], nil
}

func (p *parser) findStdRedirect() int {
	for i, s := range p.state.parts {
		if p.isStdoutRedirect(s) {
			return i
		}
	}

	return -1
}

func (p *parser) isRedirect(s string) bool {
	return p.isStdoutRedirect(s) || p.isStderrRedirect(s)
}

func (p *parser) isStdoutRedirect(s string) bool {
	return s == ">" || s == "1>"
}

func (p *parser) isStderrRedirect(s string) bool {
	return s == "2>"
}

func (p *parser) readAt(i int) (string, int) {
	i = p.skipWhitespaceFrom(i)
	token, i := p.readUntilWhitespace(i)

	return token, i
}

func (p *parser) readUntilWhitespace(i int) (string, int) {
	token := strings.Builder{}

	c, size := utf8.DecodeRuneInString(p.cmd[i:])
	for i < len(p.cmd) && !unicode.IsSpace(c) {
		if c == '\'' {
			i = p.readSingleQuotesAt(i, &token)
		} else if c == '"' {
			i = p.readDoubleQuotesAt(i, &token)
		} else if c == '\\' {
			i = p.readEscapedCharAt(i, &token)
		} else if c == '>' || c == '1' {
			token, i = p.readStdRedirect(i, &token)
		} else if c == '2' {
			token, i = p.readErrRedirect(i, &token)
		} else {
			token.WriteRune(c)
			i += size
		}

		c, size = utf8.DecodeRuneInString(p.cmd[i:])
	}

	return token.String(), i
}

func (p *parser) readErrRedirect(i int, token *strings.Builder) (strings.Builder, int) {
	c, _ := utf8.DecodeRuneInString(p.cmd[i:])
	if i < len(p.cmd)-1 && p.cmd[i+1] == '>' && token.Len() == 0 {
		if token.Len() > 0 {
			p.state.parts = append(p.state.parts, token.String())
		}
		p.state.parts = append(p.state.parts, "2>")
		return strings.Builder{}, i + 2
	} else {
		token.WriteRune(c)
		return *token, i + 1
	}
}

func (p *parser) readStdRedirect(i int, token *strings.Builder) (strings.Builder, int) {
	c, _ := utf8.DecodeRuneInString(p.cmd[i:])
	if c == '>' {
		if token.Len() > 0 {
			p.state.parts = append(p.state.parts, token.String())
		}

		if i < len(p.cmd)-1 && p.cmd[i+1] == '>' {
			p.state.parts = append(p.state.parts, ">>")
			return strings.Builder{}, i + 2
		} else {
			p.state.parts = append(p.state.parts, ">")
			return strings.Builder{}, i + 1
		}
	}

	if c != '1' {
		return *token, i
	}

	if i < len(p.cmd)-1 && p.cmd[i+1] == '>' && token.Len() == 0 {
		if token.Len() > 0 {
			p.state.parts = append(p.state.parts, token.String())
		}

		if i < len(p.cmd)-2 && p.cmd[i+2] == '>' {
			p.state.parts = append(p.state.parts, "1>>")
			return strings.Builder{}, i + 3
		} else {
			p.state.parts = append(p.state.parts, "1>")
			return strings.Builder{}, i + 2
		}
	}

	token.WriteRune(c)
	return *token, i + 1
}

func (p *parser) readEscapedCharAt(i int, token *strings.Builder) int {
	if i == len(p.cmd)-1 {
		return len(p.cmd)
	}

	i++
	c, size := utf8.DecodeRuneInString(p.cmd[i:])
	token.WriteRune(c)

	return i + size
}

func (p *parser) readDoubleQuotesAt(i int, token *strings.Builder) int {
	i++
	c, size := utf8.DecodeRuneInString(p.cmd[i:])
	for i < len(p.cmd) {
		if c == '"' {
			return i + 1
		} else if c == '\\' && p.isEscapedInDoubleQuotes(i+1) {
			i = p.readEscapedCharAt(i, token)
		} else {
			token.WriteRune(c)
			i += size
		}

		c, size = utf8.DecodeRuneInString(p.cmd[i:])
	}

	return i
}

func (p *parser) isEscapedInDoubleQuotes(i int) bool {
	if i >= len(p.cmd) {
		return false
	}

	c, _ := utf8.DecodeRuneInString(p.cmd[i:])
	return c == '"' || c == '\\' || c == '$' || c == '`'
}

func (p *parser) readSingleQuotesAt(i int, token *strings.Builder) int {
	for j, c := range p.cmd[i+1:] {
		if c == '\'' {
			return i + j + 2
		}

		token.WriteRune(c)
	}

	return len(p.cmd)
}

func (p *parser) skipWhitespaceFrom(i int) int {
	c, size := utf8.DecodeRuneInString(p.cmd[i:])
	for i < len(p.cmd) && unicode.IsSpace(c) {
		i += size
		c, size = utf8.DecodeRuneInString(p.cmd[i:])
	}

	return i
}
