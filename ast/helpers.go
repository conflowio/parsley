// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast

import "github.com/conflowio/parsley/parsley"

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

// ReaderPosSetter allows to change the reader position on a node
type ReaderPosSetter interface {
	SetReaderPos(f func(parsley.Pos) parsley.Pos)
}

// ReaderPosSetterNode defines a node which also implements the ReaderPosSetter interface
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . ReaderPosSetterNode
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
