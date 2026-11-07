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
	ok          bool
	verb        string
	params      []string
	stdRedirect string
	errRedirect string
}

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

	params := p.state.parts[1:]
	stdRedirectIdx, stdRedirect, err := p.getStdRedirect()
	if err != nil {
		fmt.Println(err.Error())
		return parseResult{
			ok: false,
		}
	} else if stdRedirectIdx > 0 {
		params = p.state.parts[1:stdRedirectIdx]
	}

	errRedirectIdx, errRedirect, err := p.getErrRedirect()
	if err != nil {
		fmt.Println(err.Error())
		return parseResult{
			ok: false,
		}
	} else if errRedirectIdx > 0 {
		params = p.state.parts[1:errRedirectIdx]
	}

	return parseResult{
		ok:          true,
		verb:        p.state.parts[0],
		params:      params,
		stdRedirect: stdRedirect,
		errRedirect: errRedirect,
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
		p.state.parts = append(p.state.parts, ">")
		return strings.Builder{}, i + 1
	}

	if c == '1' {
		if i < len(p.cmd)-1 && p.cmd[i+1] == '>' && token.Len() == 0 {
			if token.Len() > 0 {
				p.state.parts = append(p.state.parts, token.String())
			}
			p.state.parts = append(p.state.parts, "1>")
			return strings.Builder{}, i + 2
		} else {
			token.WriteRune(c)
			return *token, i + 1
		}
	}

	return *token, i
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
