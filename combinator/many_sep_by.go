package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/parser"
)

// ManySepBy matches the given value parser one or more times separated by the separator parser
func ManySepBy(name string, token string, h *parser.History, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Func {
	sepValue := Memoize(name+"_SV", h, And(builder.All("TMP", interpreter), sepP, valueP))
	sepValueMany := Memoize(name+"_SV*", h, Many(builder.Flatten(token, interpreter), sepValue))
	return And(mergeChildren(token, interpreter), valueP, Maybe(sepValueMany))
}

func mergeChildren(token string, interpreter ast.Interpreter) ast.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		children := []ast.Node{nodes[0]}
		if len(nodes) > 1 {
			node1 := nodes[1].(ast.NonTerminalNode)
			children = append(children, node1.Children()...)
		}
		return ast.NewNonTerminalNode(token, children, interpreter)
	})
}
