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

// WalkNode applies the given function to the node
func WalkNode(node parsley.Node, f func(n parsley.Node) bool) bool {
	switch n := node.(type) {
	case Walkable:
		return n.Walk(f)
	default:
		return f(node)
	}
}

// ReaderPosSetter allows to change the reader position on a node
type ReaderPosSetter interface {
	SetReaderPos(f func(parsley.Pos) parsley.Pos)
}

// SetReaderPos sets the reader position on a node
func SetReaderPos(node parsley.Node, f func(parsley.Pos) parsley.Pos) parsley.Node {
	switch n := node.(type) {
	case ReaderPosSetter:
		n.SetReaderPos(f)
	case NilNode:
		return NilNode(f(parsley.Pos(n)))
	case NodeList:
		for i, item := range n {
			n[i] = SetReaderPos(item, f)
		}
	default:
		panic("invalid node type for SetReaderPos(), you may need to implement the ast.ReaderPosSetter interface")
	}
	return node
}
