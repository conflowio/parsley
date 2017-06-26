package terminal

import (
	"strconv"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/token"
)

// Float matches a float literal
func Float() parser.Func {
	return parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		tr := r.(*text.Reader)
		cur := tr.Cursor()
		if matches, pos, ok := tr.ReadMatch("[-+]?[0-9]*\\.[0-9]+(?:[eE][-+]?[0-9]+)?", false); ok {
			val, err := strconv.ParseFloat(matches[0], 64)
			if err != nil {
				return parser.NoCurtailingParsers(), nil, parser.NewError(cur, "invalid float value encountered")
			}
			var rs parser.ResultSet
			rs = parser.NewResult(ast.NewTerminalNode(token.FLOAT, pos, val), r).AsSet()
			return parser.NoCurtailingParsers(), rs, nil
		}
		return parser.NoCurtailingParsers(), nil, parser.NewError(cur, "was expecting float value")
	})
}
