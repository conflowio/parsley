// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package builder defines functions for generating basic node builder functions
package builder

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"
)

// All returns with a node builder function for including all nodes for interpretation
func All(token string, interpreter parsley.Interpreter) parsley.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []parsley.Node) parsley.Node {
		return ast.NewNonTerminalNode(token, nodes).Bind(interpreter)
	})
}

// BinaryOperation returns with a node builder function for building binary operator nodes
func BinaryOperation(interpreter parsley.Interpreter) parsley.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []parsley.Node) parsley.Node {
		if len(nodes) != 3 {
			panic("BinaryOperation builder should receive exactly three nodes")
		}
		return ast.NewNonTerminalNode(nodes[1].Token(), []parsley.Node{nodes[0], nodes[2]}).Bind(interpreter)
	})
}

// Flatten returns with a node builder function which puts all nodes and their direct children
// flattened in a new node
func Flatten(token string, interpreter parsley.Interpreter) parsley.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []parsley.Node) parsley.Node {
		var children []parsley.Node
		for _, node := range nodes {
			switch n := node.(type) {
			case *ast.TerminalNode:
				children = append(children, n)
			case *ast.NonTerminalNode:
				children = append(children, n.Children()...)
			}
		}
		return ast.NewNonTerminalNode(token, children).Bind(interpreter)
	})
}
