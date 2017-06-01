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

// Char matches a character literal enclosed in single quotes
func Char() parser.Func {
	return parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		tr := r.(*text.Reader)
		if matches, pos := tr.ReadMatch("^'(.|\\\\[abfnrtv]|\\\\x[0-9a-fA-F]{2,2}|\\\\u[0-9a-fA-F]{4,4}|\\\\U[0-9a-fA-F]{8,8})'"); matches != nil {
			match := matches[1]
			value, _, tail, err := strconv.UnquoteChar(match, '"')
			if tail != "" {
				panic(fmt.Sprintf("Unprocessed string segment: %s", tail))
			}
			if err == nil {
				return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(token.CHAR, pos, value), r).AsSet()
			}
		}
		return parser.NoCurtailingParsers(), nil
	})
}
