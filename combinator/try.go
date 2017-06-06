package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Try combines the the parsers until it fails and returns the longest match
func Try(nodeBuilder ast.NodeBuilder, min int, parsers ...parser.Parser) parser.Func {
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
		return NewRecursive(nodeBuilder, lookup, min, len(parsers)).Parse(leftRecCtx, r)
	})
}
