package terminal

import (
	"regexp"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
)

// Rune matches the given character
func Rune(char rune, token string) parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		tr := r.(*text.Reader)
		if _, pos, ok := tr.ReadMatch(regexp.QuoteMeta(string(char))); ok {
			return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(token, pos, char), r).AsSet()
		}
		return parser.NoCurtailingParsers(), nil
	})
}
