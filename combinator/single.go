package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// Single will change the result of p if it returns with a non terminal node
// with only one child. In this case directly the child will returned.
func Single(p parsley.Parser) parser.Func {
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		res, cp, err := p.Parse(ctx, leftRecCtx, pos)
		if err != nil {
			return nil, cp, err
		}

		if node, ok := res.(*ast.NonTerminalNode); ok {
			if len(node.Children()) == 1 {
				return node.Children()[0], cp, nil
			}
		}

		return res, cp, nil
	})
}
