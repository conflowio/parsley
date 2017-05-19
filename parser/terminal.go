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
	return Func(func(h *History, r *reader.Reader) Results {
		if r.IsEOF() {
			return NewResults(Result{ast.NewTerminalNode(token.EOF, r.Cursor(), nil), r})
		}
		return nil
	})
}

// Rune matches one specific character
func Rune(char rune, token token.Token) Func {
	return Func(func(h *History, r *reader.Reader) Results {
		if matches, pos := r.ReadMatch("^\\s?" + regexp.QuoteMeta(string(char)) + "\\s?"); matches != nil {
			return NewResults(Result{ast.NewTerminalNode(token, pos, string(char)), r})
		}
		return nil
	})
}

// CharLiteral matches a character literal enclosed in single quotes
func CharLiteral() Func {
	return Func(func(h *History, r *reader.Reader) Results {
		if matches, pos := r.ReadMatch("^'(.)'"); matches != nil {
			return NewResults(Result{ast.NewTerminalNode(token.CHAR, pos, matches[1]), r})
		}
		return nil
	})
}

// IntLiteral matches an integer literal
func IntLiteral() Func {
	return Func(func(h *History, r *reader.Reader) Results {
		if matches, pos := r.ReadMatch("^[\\-+]?[1-9][0-9]*"); matches != nil {
			intValue, err := strconv.Atoi(matches[0])
			if err != nil {
				panic(fmt.Sprintf("Could not convert %s to integer", matches[0]))
			}
			return NewResults(Result{ast.NewTerminalNode(token.INT, pos, intValue), r})
		}
		return nil
	})
}
