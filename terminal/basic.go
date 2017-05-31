package terminal

import (
	"regexp"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Empty always matches and returns with an empty node result
func Empty() parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		return parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(ast.EMPTY, r.Cursor(), nil), r))
	})
}

// End matches the end of the input
func End() parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		if r.IsEOF() {
			return parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(ast.EOF, r.Cursor(), nil), r))
		}
		return nil
	})
}

// Rune matches one specific character
func Rune(char rune, token string) parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		if matches, pos := r.ReadMatch("^" + regexp.QuoteMeta(string(char))); matches != nil {
			return parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(token, pos, char), r))
		}
		return nil
	})
}
