package main

import (
	"fmt"
	"strings"
	"unicode"
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
	for _, c := range p.cmd {
		p.processChar(c)
	}

	if p.state.currentWord.Len() > 0 {
		p.state.parts = append(p.state.parts, p.state.currentWord.String())
		p.state.currentWord.Reset()
	}

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
		return
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
