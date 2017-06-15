package terminal

import (
	"regexp"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
)

// Word matches the given word
func Word(word string, token string, value interface{}) parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		tr := r.(*text.Reader)
		if _, pos, ok := tr.ReadMatch(regexp.QuoteMeta(word), false); ok {
			if _, ok := tr.PeakMatch("\\w+"); ok {
				return parser.NoCurtailingParsers(), nil
			}
			return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(token, pos, value), r).AsSet()
		}
		return parser.NoCurtailingParsers(), nil
	})
}
