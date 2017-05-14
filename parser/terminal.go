package parser

import (
	"fmt"
	"go/token"
	"regexp"
	"strconv"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/reader"
)

// End matches the end of the input
func End() Func {
	return Func(func(r *reader.Reader) (ast.Node, *reader.Reader) {
		if r.IsEOF() {
			return ast.NewTerminalNode(token.EOF, ""), r
		}
		return nil, nil
	})
}

// Rune matches one specific character
func Rune(char rune, token token.Token) Func {
	return Func(func(r *reader.Reader) (ast.Node, *reader.Reader) {
		if matches := r.ReadMatch("^\\s?" + regexp.QuoteMeta(string(char)) + "\\s?"); matches != nil {
			return ast.NewTerminalNode(token, string(char)), r
		}
		return nil, nil
	})
}

// CharLiteral matches a character literal enclosed in single quotes
func CharLiteral() Func {
	return Func(func(r *reader.Reader) (ast.Node, *reader.Reader) {
		if matches := r.ReadMatch("^'(.)'"); matches != nil {
			return ast.NewTerminalNode(token.CHAR, matches[1]), r
		}
		return nil, nil
	})
}

// IntLiteral matches an integer literal
func IntLiteral() Func {
	return Func(func(r *reader.Reader) (ast.Node, *reader.Reader) {
		if matches := r.ReadMatch("^[\\-+]?[1-9][0-9]*"); matches != nil {
			intValue, err := strconv.Atoi(matches[0])
			if err != nil {
				panic(fmt.Sprintf("Could not convert %s to integer", matches[0]))
			}
			return ast.NewTerminalNode(token.INT, intValue), r
		}
		return nil, nil
	})
}
