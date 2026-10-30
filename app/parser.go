package main

import (
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

	parts       []string
	currentWord strings.Builder
}

func buildParser(cmd string) parser {
	return parser{
		cmd: cmd,
	}
}

func (p *parser) parse() []string {
	for _, c := range p.cmd {
		p.processChar(c)
	}

	if p.state.currentWord.Len() > 0 {
		p.state.parts = append(p.state.parts, p.state.currentWord.String())
		p.state.currentWord.Reset()
	}

	return p.state.parts
}

func (p *parser) processChar(c rune) {
	if c == '\\' {
		if !p.isEscaped() {
			p.state.isBackslashEscaped = true
			return
		} else if p.state.isDoubleQuotesEscaped {
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
