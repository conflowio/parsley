package ast

import (
	"fmt"
)

// NodeBuilder defines an interface to build a node from multiple nodes
type NodeBuilder interface {
	BuildNode([]Node) Node
}

// NodeBuilderFunc defines a helper to implement the NodeBuilder interface with functions
type NodeBuilderFunc func([]Node) Node

// BuildNode combines the given nodes into a single node
func (f NodeBuilderFunc) BuildNode(nodes []Node) Node {
	return f(nodes)
}

// SingleNodeBuilder simply returns with the node at the given index
func SingleNodeBuilder(i int) NodeBuilder {
	return NodeBuilderFunc(func(nodes []Node) Node {
		if i >= len(nodes) {
			panic(fmt.Sprintf("Node index is out of bounds: %d", i))
		}
		return nodes[i]
	})
}

// AllNodesBuilder returns with a function for including all nodes for interpretation
func AllNodesBuilder(token string, interpreter Interpreter) NodeBuilder {
	return NodeBuilderFunc(func(nodes []Node) Node {
		return NewNonTerminalNode(
			token,
			nodes,
			interpreter,
		)
	})
}

// BinaryOperatorBuilder returns with a function for building binary operator nodes
func BinaryOperatorBuilder(interpreter Interpreter) NodeBuilder {
	return NodeBuilderFunc(func(nodes []Node) Node {
		if len(nodes) != 3 {
			panic("BinaryOperatorBuilder should receive exactly three nodes")
		}
		return NewNonTerminalNode(
			nodes[1].Token(),
			[]Node{nodes[0], nodes[2]},
			interpreter,
		)
	})
}

// NilBuilder returns with a node builder function which always returns with nil
func NilBuilder() NodeBuilder {
	return NodeBuilderFunc(func(nodes []Node) Node {
		return nil
	})
}

// FlattenBuilder returns all nodes and their direct children flattened in a new node
func FlattenBuilder(token string, interpreter Interpreter) NodeBuilder {
	return NodeBuilderFunc(func(nodes []Node) Node {
		var children []Node
		for _, node := range nodes {
			switch n := node.(type) {
			case TerminalNode:
				children = append(children, n)
			case NonTerminalNode:
				children = append(children, n.children...)
			}
		}
		return NewNonTerminalNode(token, children, interpreter)
	})
}
