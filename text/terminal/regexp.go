package terminal

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
)

// Regexp matches the given regular expression
func Regexp(regexp string, groupIndex int, token string) parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		tr := r.(*text.Reader)
		if matches, pos, ok := tr.ReadMatch(regexp, false); ok {
			if groupIndex >= len(matches) {
				panic(fmt.Sprintf("Capturing group %d is invalid for %s", groupIndex, regexp))
			}
			if token == "" {
				token = matches[groupIndex]
			}
			return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(token, pos, matches[groupIndex]), r).AsSet()
		}
		return parser.NoCurtailingParsers(), nil
	})
}
