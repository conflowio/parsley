package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// And applies the given parsers in order and returns the result only if all parsers matched
func And(nodeBuilder ast.NodeBuilder, parsers ...parser.Parser) parser.Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	lookup := func(i int) parser.Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		l := len(parsers)
		return newRecursive(nodeBuilder, lookup, l, l).Parse(leftRecCtx, r)
	})
}
