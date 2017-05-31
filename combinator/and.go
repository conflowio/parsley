package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// And combines multiple parsers
func And(nodeBuilder ast.NodeBuilder, parsers ...parser.Parser) parser.Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	return parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		return NewRecursive(nodeBuilder, false, parserListLookUp(parsers)).Parse(leftRecCtx, r)
	})
}

func parserListLookUp(parsers []parser.Parser) func(i int) parser.Parser {
	return func(i int) parser.Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
}
