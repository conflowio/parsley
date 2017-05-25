package ast

import (
	"fmt"
)

// NodeBuilder defines a function to build a node from multiple nodes
type NodeBuilder func([]Node) Node

// BinaryOperatorBuilder returns with a function for building binary operator nodes
func BinaryOperatorBuilder(token string, interpreter Interpreter) NodeBuilder {
	return func(nodes []Node) Node {
		if len(nodes) != 3 {
			panic("BinaryOperatorBuilder should receive exactly three nodes")
		}
		return NewNonTerminalNode(
			token,
			[]Node{nodes[0], nodes[2]},
			interpreter,
		)
	}
}

// SingleNodeBuilder simply returns with the node at the given index
func SingleNodeBuilder(i int) NodeBuilder {
	return func(nodes []Node) Node {
		if i >= len(nodes) {
			panic(fmt.Sprintf("Node index is out of bounds: %d", i))
		}
		return nodes[i]
	}
}
