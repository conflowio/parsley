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

// Bool matches a bool literal: true or false
func Bool() parser.Func {
	return parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		tr := r.(*text.Reader)
		if matches, pos, ok := tr.ReadMatch("true|false", false); ok {
			if _, ok := tr.PeakMatch("\\w+"); ok {
				return parser.NoCurtailingParsers(), nil
			}
			val, err := strconv.ParseBool(matches[0])
			if err != nil {
				panic(fmt.Sprintf("Invalid bool value encountered: %s", matches[0]))
			}
			return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(token.BOOL, pos, val), r).AsSet()
		}
		return parser.NoCurtailingParsers(), nil
	})
}
