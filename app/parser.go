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
	isSingleQuotesEscaped bool
	isDoubleQuotesEscaped bool
	isBackslashEscaped    bool

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

	// fmt.Println("parts:", p.state.parts)
	return p.state.parts
	// var res []string
	// var currentWord strings.Builder
	//
	// for _, char := range p.cmd {
	// 	if char == '\\' && !isSingleQuotesEscaped && !isDoubleQuotesEscaped && !isBackslashEscaped {
	// 		isBackslashEscaped = true
	// 		continue
	// 	}
	//
	// 	if char == '"' && !isSingleQuotesEscaped && !isBackslashEscaped {
	// 		isDoubleQuotesEscaped = !isDoubleQuotesEscaped
	// 		continue
	// 	}
	//
	// 	if char == '\'' && !isDoubleQuotesEscaped && !isBackslashEscaped {
	// 		isSingleQuotesEscaped = !isSingleQuotesEscaped
	// 		continue
	// 	}
	//
	// 	if unicode.IsSpace(char) {
	// 		if isSingleQuotesEscaped || isDoubleQuotesEscaped || isBackslashEscaped {
	// 			currentWord.WriteRune(char)
	// 		} else if currentWord.Len() > 0 {
	// 			res = append(res, currentWord.String())
	// 			currentWord.Reset()
	// 		}
	// 	} else {
	// 		currentWord.WriteRune(char)
	// 	}
	//
	// 	if isBackslashEscaped {
	// 		isBackslashEscaped = false
	// 	}
	// }
	//
	// if currentWord.Len() > 0 {
	// 	res = append(res, currentWord.String())
	// }
	//
	// return res
}

func (p *parser) processChar(c rune) {
	// fmt.Println("char:", string(c))
	// fmt.Println("state:", p.state)
	if c == '\\' && !p.isEscaped() {
		p.state.isBackslashEscaped = true
		return
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
			// fmt.Println("word done")
			p.state.parts = append(p.state.parts, p.state.currentWord.String())
			p.state.currentWord.Reset()
		}
	} else {
		p.state.currentWord.WriteRune(c)
		// fmt.Println("writing letter")
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
