// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast

import (
	"fmt"

	"github.com/opsidian/parsley/parsley"
)

// NIL is the nil token
const NIL = "NIL"

// EOF is the end of file token
const EOF = "EOF"

// TerminalNode is a leaf node in the AST
type TerminalNode struct {
	token     string
	value     interface{}
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewTerminalNode creates a new TerminalNode instance
func NewTerminalNode(token string, value interface{}, pos parsley.Pos, readerPos parsley.Pos) *TerminalNode {
	return &TerminalNode{
		token:     token,
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (t *TerminalNode) Token() string {
	return t.token
}

// Value returns with the value of the node
func (t *TerminalNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	return t.value, nil
}

// Pos returns the position
func (t *TerminalNode) Pos() parsley.Pos {
	return t.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (t *TerminalNode) ReaderPos() parsley.Pos {
	return t.readerPos
}

// SetReaderPos changes the reader position
func (t *TerminalNode) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	t.readerPos = f(t.readerPos)
}

// String returns with a string representation of the node
func (t *TerminalNode) String() string {
	return fmt.Sprintf("%s{%v, %d..%d}", t.token, t.value, t.pos, t.readerPos)
}

// NilNode represents an nil node
type NilNode parsley.Pos

// Token returns with NIL
func (n NilNode) Token() string {
	return NIL
}

// Value returns with nil
func (n NilNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	return nil, nil
}

// Pos returns with NilPosition
func (n NilNode) Pos() parsley.Pos {
	return parsley.Pos(n)
}

// ReaderPos returns the reader position
func (n NilNode) ReaderPos() parsley.Pos {
	return parsley.Pos(n)
}

// String returns with a string representation of the node
func (n NilNode) String() string {
	return NIL
}

// NonTerminalNode represents a branch node in the AST
type NonTerminalNode struct {
	token       string
	children    []parsley.Node
	pos         parsley.Pos
	readerPos   parsley.Pos
	interpreter parsley.Interpreter
}

// NewNonTerminalNode creates a new NonTerminalNode instance
func NewNonTerminalNode(token string, children []parsley.Node, interpreter parsley.Interpreter) *NonTerminalNode {
	if len(children) == 0 {
		panic("NewNonTerminalNode should not be called with empty node list")
	}
	for _, c := range children {
		if c == nil {
			panic("NewNonTerminalNode can not have children with nil values")
		}
	}
	return &NonTerminalNode{
		token:       token,
		children:    children,
		pos:         children[0].Pos(),
		readerPos:   children[len(children)-1].ReaderPos(),
		interpreter: interpreter,
	}
}

// NewEmptyNonTerminalNode creates a new NonTerminalNode without children
func NewEmptyNonTerminalNode(token string, pos parsley.Pos, interpreter parsley.Interpreter) *NonTerminalNode {
	return &NonTerminalNode{
		token:       token,
		children:    []parsley.Node{},
		pos:         pos,
		readerPos:   pos,
		interpreter: interpreter,
	}
}

// Token returns with the node token
func (n *NonTerminalNode) Token() string {
	return n.token
}

// Value returns with the value of the node
func (n *NonTerminalNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	if n.interpreter == nil {
		panic("missing interpreter for node")
	}
	return n.interpreter.Eval(ctx, n.children)
}

// Pos returns the position
func (n *NonTerminalNode) Pos() parsley.Pos {
	return n.pos
}

// Children returns with the children
func (n *NonTerminalNode) Children() []parsley.Node {
	return n.children
}

// ReaderPos returns the position of the first character immediately after this node
func (n *NonTerminalNode) ReaderPos() parsley.Pos {
	return n.readerPos
}

// SetReaderPos amends the reader position using the given function
func (n *NonTerminalNode) SetReaderPos(f func(parsley.Pos) parsley.Pos) {
	n.readerPos = f(n.readerPos)
}

// String returns with a string representation of the node
func (n *NonTerminalNode) String() string {
	return fmt.Sprintf("%s{%s, %d..%d}", n.token, n.children, n.pos, n.readerPos)
}

// Walk runs the given function on all child nodes
func (n *NonTerminalNode) Walk(f func(i int, n parsley.Node) bool) bool {
	for i, node := range n.children {
		if f(i, node) {
			return true
		}
	}

	return false
}

// NodeList contains a list of nodes, should be used when a parser returns with multiple results
type NodeList []parsley.Node

// Token returns with NODE_LIST
func (nl NodeList) Token() string {
	return "NODE_LIST"
}

// Value returns with the value of the first result
func (nl NodeList) Value(ctx interface{}) (interface{}, parsley.Error) {
	if len(nl) > 0 {
		return nl[0].Value(ctx)
	}
	return nil, nil
}

// Pos returns the value of the first pos (all nodes should have the same position)
func (nl NodeList) Pos() parsley.Pos {
	if len(nl) > 0 {
		return nl[0].Pos()
	}
	return parsley.NilPos
}

// ReaderPos should not be called on a NodeList
func (nl NodeList) ReaderPos() parsley.Pos {
	panic("ReaderPos() should not be called on NodeList")
}

// Append appends a new node to the list
func (nl *NodeList) Append(node parsley.Node) {
	switch v := node.(type) {
	case NodeList:
		for _, node := range v {
			nl.Append(node)
		}
	case NilNode:
		for _, node := range *nl {
			if node == v {
				return
			}
		}
		*nl = append(*nl, v)
	default:
		*nl = append(*nl, v)
	}
}

// Walk runs the given function on all nodes
func (nl NodeList) Walk(f func(i int, n parsley.Node) bool) bool {
	for i, node := range nl {
		if f(i, node) {
			return true
		}
	}

	return false
}
