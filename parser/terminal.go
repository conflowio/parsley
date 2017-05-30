package parser

import (
	"regexp"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/reader"
)

// Empty always matches and returns with an empty node result
func Empty() Func {
	return Func(func(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
		return NewParserResult(NoCurtailingParsers(), NewResult(ast.NewTerminalNode(ast.EMPTY, r.Cursor(), nil), r))
	})
}

// End matches the end of the input
func End() Func {
	return Func(func(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
		if r.IsEOF() {
			return NewParserResult(NoCurtailingParsers(), NewResult(ast.NewTerminalNode(ast.EOF, r.Cursor(), nil), r))
		}
		return nil
	})
}

// Rune matches one specific character
func Rune(char rune, token string) Func {
	return Func(func(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
		if matches, pos := r.ReadMatch("^" + regexp.QuoteMeta(string(char))); matches != nil {
			return NewParserResult(NoCurtailingParsers(), NewResult(ast.NewTerminalNode(token, pos, char), r))
		}
		return nil
	})
}
