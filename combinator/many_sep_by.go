package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/parser"
)

// ManySepBy matches the given value parser zero or more times separated by the separator parser
func ManySepBy(token string, h *parser.History, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Func {
	sepValue := Memoize(token+"_MSB", h, And(builder.All("SEP_VALUE", nil), sepP, valueP))
	sepValueMany := Memoize(token+"_MSB*", h, Many(builder.Flatten(token, nil), sepValue, 0, -1))
	return Try(mergeChildren(token, interpreter), 0, valueP, sepValueMany)
}

// ManySepBy1 matches the given value parser one or more times separated by the separator parser
func ManySepBy1(token string, h *parser.History, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Func {
	sepValue := Memoize(token+"_MSB", h, And(builder.All("SEP_VALUE", nil), sepP, valueP))
	sepValueMany := Memoize(token+"_MSB*", h, Many(builder.Flatten(token, nil), sepValue, 0, -1))
	return Try(mergeChildren(token, interpreter), 1, valueP, sepValueMany)
}

func mergeChildren(token string, interpreter ast.Interpreter) ast.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		if len(nodes) == 0 {
			return ast.NewNonTerminalNode(token, nil, interpreter)
		}

		children := []ast.Node{nodes[0]}
		if len(nodes) > 1 && nodes[1] != nil {
			node1 := nodes[1].(ast.NonTerminalNode)
			children = append(children, node1.Children()...)
		}
		return ast.NewNonTerminalNode(token, children, interpreter)
	})
}
