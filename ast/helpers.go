// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast

import "github.com/opsidian/parsley/parsley"

// AppendNode appends
func AppendNode(n1, n2 parsley.Node) parsley.Node {
	if n1 == nil {
		return n2
	}
	if n2 == nil {
		return n1
	}
	switch n := n1.(type) {
	case NodeList:
		n.Append(n2)
		return n
	default:
		nl := NodeList([]parsley.Node{n1})
		nl.Append(n2)
		return nl
	}
}

// Walkable is a generic interface to allow to apply a function on the node
// The Walk function should return true if the walk should be interrupted
type Walkable interface {
	Walk(f func(n parsley.Node) bool) bool
}

// WalkableNode defines a node which also implements the Walkable interface
//go:generate counterfeiter . WalkableNode
type WalkableNode interface {
	parsley.Node
	Walk(f func(n parsley.Node) bool) bool
}

// WalkNode applies the given function to the node
func WalkNode(node parsley.Node, f func(n parsley.Node) bool) bool {
	switch n := node.(type) {
	case Walkable:
		return n.Walk(f)
	case parsley.NonTerminalNode:
		for _, child := range n.Children() {
			if WalkNode(child, f) {
				return true
			}
		}
		return false
	default:
		return f(node)
	}
}

// ReaderPosSetter allows to change the reader position on a node
type ReaderPosSetter interface {
	SetReaderPos(f func(parsley.Pos) parsley.Pos)
}

// ReaderPosSetterNode defines a node which also implements the ReaderPosSetter interface
//go:generate counterfeiter . ReaderPosSetterNode
type ReaderPosSetterNode interface {
	parsley.Node
	ReaderPosSetter
}

// SetReaderPos sets the reader position on a node
func SetReaderPos(node parsley.Node, f func(parsley.Pos) parsley.Pos) parsley.Node {
	switch n := node.(type) {
	case ReaderPosSetter:
		n.SetReaderPos(f)
	case EmptyNode:
		return EmptyNode(f(parsley.Pos(n)))
	default:
		panic("invalid node type for SetReaderPos(), you need to implement the ast.ReaderPosSetter interface")
	}
	return node
}
