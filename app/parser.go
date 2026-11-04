package main

import (
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
	isSingleQuotesEscaped            bool
	isDoubleQuotesEscaped            bool
	isBackslashEscaped               bool
	isBackslashInDoubleQuotesEscaped bool
	isStdRedirectMaybe               bool
	isErrRedirectMaybe               bool

	parts       []string
	currentWord strings.Builder
}

type parseResult struct {
	ok          bool
	verb        string
	params      []string
	stdRedirect string
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

	// fmt.Println(p.state.parts)
	// for _, c := range p.cmd {
	// 	p.processChar(c)
	// }

	// todo: tihs needs to go
	// if p.state.isStdRedirectMaybe {
	// 	p.state.currentWord.WriteRune('1')
	// }

	// if p.state.currentWord.Len() > 0 {
	// 	p.state.parts = append(p.state.parts, p.state.currentWord.String())
	// 	p.state.currentWord.Reset()
	// }

	if len(p.state.parts) == 0 {
		return parseResult{
			ok: true,
		}
	}

	params := p.state.parts[1:]
	stdRedirect := ""
	stdRedirectIdx := p.findStdRedirect()
	if stdRedirectIdx >= 0 {
		if len(p.state.parts) <= stdRedirectIdx+1 {
			fmt.Println("nsh: syntax error near unexpected token `newline'")
			return parseResult{
				ok: false,
			}
		}

		stdRedirect = p.state.parts[stdRedirectIdx+1]
		params = p.state.parts[1:stdRedirectIdx]
	}

	return parseResult{
		ok:          true,
		verb:        p.state.parts[0],
		params:      params,
		stdRedirect: stdRedirect,
	}
}

func (p *parser) readAt(i int) (string, int) {
	i = p.skipWhitespaceFrom(i)

	token, i := p.readUntilWhitespace(i)
	// todo: redirect
	//

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
		} else {
			token.WriteRune(c)
			i += size
		}

		c, size = utf8.DecodeRuneInString(p.cmd[i:])
	}

	return token.String(), i
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
			// fmt.Println("returning:", i+1)
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
			// fmt.Println("returning:", i+j+2)
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

func (p *parser) processChar(c rune) {
	if c == '\\' {
		if !p.isEscaped() {
			p.state.isBackslashEscaped = true
			return
		} else if p.state.isDoubleQuotesEscaped && !p.state.isBackslashInDoubleQuotesEscaped {
			p.state.isBackslashInDoubleQuotesEscaped = true
			return
		}
	}

	if p.state.isBackslashInDoubleQuotesEscaped {
		p.state.isBackslashInDoubleQuotesEscaped = false

		if c == '"' || c == '\\' || c == '$' || c == '`' {
			p.state.currentWord.WriteRune(c)
			return
		} else if c == '\n' {
			return
		} else {
			p.state.currentWord.WriteRune('\\')
			p.state.currentWord.WriteRune(c)
			return
		}
	}

	if c == '"' && !p.state.isSingleQuotesEscaped && !p.state.isBackslashEscaped {
		p.state.isDoubleQuotesEscaped = !p.state.isDoubleQuotesEscaped
		return
	}

	if c == '\'' && !p.state.isDoubleQuotesEscaped && !p.state.isBackslashEscaped {
		p.state.isSingleQuotesEscaped = !p.state.isSingleQuotesEscaped
		return
	}

	// todo: refactor this shit to look at multiple chars
	if c == '>' && !p.isEscaped() {
		if p.state.currentWord.Len() > 0 {
			p.state.parts = append(p.state.parts, p.state.currentWord.String())
			p.state.currentWord.Reset()
		}

		p.state.currentWord.WriteRune(c)
		p.state.parts = append(p.state.parts, p.state.currentWord.String())
		p.state.currentWord.Reset()
		p.state.isStdRedirectMaybe = false
		return
	}
	// if c == '>' && (p.state.currentWord.Len() == 0 || p.state.isStdRedirectMaybe) {
	// 	p.state.currentWord.WriteRune(c)
	// 	p.state.parts = append(p.state.parts, p.state.currentWord.String())
	// 	p.state.currentWord.Reset()
	// 	return
	// }

	if p.state.isStdRedirectMaybe {
		p.state.isStdRedirectMaybe = false
		p.state.currentWord.WriteRune('1')
	}

	if c == '1' && !p.isEscaped() {
		p.state.isStdRedirectMaybe = true
		return
	}

	if c == '2' && !p.isEscaped() {
		p.state.isErrRedirectMaybe = true
	}

	if unicode.IsSpace(c) {
		if p.isEscaped() {
			p.state.currentWord.WriteRune(c)
		} else if p.state.currentWord.Len() > 0 {
			p.state.parts = append(p.state.parts, p.state.currentWord.String())
			p.state.currentWord.Reset()
		}
	} else {
		p.state.currentWord.WriteRune(c)
	}

	if p.state.isBackslashEscaped {
		p.state.isBackslashEscaped = false
	}
}

func (p *parser) isEscaped() bool {
	return p.state.isSingleQuotesEscaped ||
		p.state.isDoubleQuotesEscaped ||
		p.state.isBackslashEscaped
}

func (p *parser) findStdRedirect() int {
	for i, s := range p.state.parts {
		if s == ">" || s == "1>" {
			return i
		}
	}

	return -1
}
