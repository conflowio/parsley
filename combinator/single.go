package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// Single will change the result of p if it returns with a non terminal node
// with only one child. In this case directly the child will returned.
func Single(p parsley.Parser) *parser.NamedFunc {
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		res, cp := p.Parse(ctx, leftRecCtx, pos)
		if node, ok := res.(*ast.NonTerminalNode); ok {
			if len(node.Children()) == 1 {
				return node.Children()[0], cp
			}
		}

		return res, cp
	}).WithName(p.Name)
}
