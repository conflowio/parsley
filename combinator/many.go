package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Many matches the same expression one or more times
func Many(nodeBuilder ast.NodeBuilder, p parser.Parser) parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		return NewRecursive(nodeBuilder, true, func(i int) parser.Parser { return p }).Parse(leftRecCtx, r)
	})
}
