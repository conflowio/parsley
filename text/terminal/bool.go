package terminal

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
)

// Bool matches a bool literal: true or false
func Bool() parser.Func {
	return parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		tr := r.(*text.Reader)
		if matches, pos, ok := tr.ReadMatch("true|false", false); ok {
			if _, ok := tr.PeakMatch("\\w+"); ok {
				return parser.NoCurtailingParsers(), nil, parser.NewError(pos, "was expecting boolean")
			}
			val := false
			if matches[0] == "true" {
				val = true
			}
			return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("BOOL", pos, val), r).AsSet(), nil
		}
		return parser.NoCurtailingParsers(), nil, parser.NewError(r.Cursor(), "was expecting boolean")
	})
}
