package builder

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
)

// Select simply returns with the node at the given index
func Select(i int) ast.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		if i >= len(nodes) {
			panic(fmt.Sprintf("Node index is out of bounds: %d", i))
		}
		return nodes[i]
	})
}

// All returns with a function for including all nodes for interpretation
func All(token string, interpreter ast.Interpreter) ast.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		return ast.NewNonTerminalNode(
			token,
			nodes,
			interpreter,
		)
	})
}

// BinaryOperation returns with a function for building binary operator nodes
func BinaryOperation(interpreter ast.Interpreter) ast.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		if len(nodes) != 3 {
			panic("BinaryOperatorBuilder should receive exactly three nodes")
		}
		return ast.NewNonTerminalNode(
			nodes[1].Token(),
			[]ast.Node{nodes[0], nodes[2]},
			interpreter,
		)
	})
}

// Nil returns with a node builder function which always returns with nil
func Nil() ast.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		return nil
	})
}

// Flatten returns all nodes and their direct children flattened in a new node
func Flatten(token string, interpreter ast.Interpreter) ast.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		var children []ast.Node
		for _, node := range nodes {
			switch n := node.(type) {
			case ast.TerminalNode:
				children = append(children, n)
			case ast.NonTerminalNode:
				children = append(children, n.Children()...)
			}
		}
		return ast.NewNonTerminalNode(token, children, interpreter)
	})
}
