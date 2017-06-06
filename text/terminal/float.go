package terminal

import (
	"fmt"
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
	return parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		tr := r.(*text.Reader)
		if matches, pos := tr.ReadMatch("[-+]?[0-9]*\\.[0-9]+(?:[eE][-+]?[0-9]+)?"); matches != nil {
			val, err := strconv.ParseFloat(matches[0], 64)
			if err != nil {
				panic(fmt.Sprintf("Invalid float value encountered: %s", matches[0]))
			}
			var rs parser.ResultSet
			rs = parser.NewResult(ast.NewTerminalNode(token.FLOAT, pos, val), r).AsSet()
			return parser.NoCurtailingParsers(), rs
		}
		return parser.NoCurtailingParsers(), nil
	})
}
