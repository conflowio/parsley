package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Seq tries to apply all parsers after each other matching effectively a sequence of tokens and returns with all // combination of the results
func Seq(nodeBuilder ast.NodeBuilder, parsers ...parser.Parser) parser.Func {
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
