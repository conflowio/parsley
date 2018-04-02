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

// EOF is the end of file token
const EOF = "EOF"

// TerminalNode is a leaf node in the AST
type TerminalNode struct {
	token     string
	value     interface{}
	pos       parsley.Pos
	readerPos int
}

// NewTerminalNode creates a new TerminalNode instance
func NewTerminalNode(token string, value interface{}, pos parsley.Pos, readerPos int) *TerminalNode {
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
func (t *TerminalNode) ReaderPos() int {
	return t.readerPos
}

// String returns with a string representation of the node
func (t *TerminalNode) String() string {
	if t.value == nil {
		return fmt.Sprintf("T{%s, %d, %d}", t.token, t.pos, t.readerPos)
	}
	return fmt.Sprintf("T{%v, %d, %d}", t.value, t.pos, t.readerPos)
}

// NonTerminalNode represents a branch node in the AST
type NonTerminalNode struct {
	token       string
	children    []parsley.Node
	pos         parsley.Pos
	readerPos   int
	interpreter parsley.Interpreter
}

// NewNonTerminalNode creates a new NonTerminalNode instance
func NewNonTerminalNode(token string, children []parsley.Node) *NonTerminalNode {
	node := &NonTerminalNode{
		token:    token,
		children: children,
	}

	for _, child := range children {
		if child != nil {
			node.pos = child.Pos()
			break
		}
	}

	for i := len(children) - 1; i >= 0; i-- {
		if children[i] != nil {
			node.readerPos = children[i].ReaderPos()
			break
		}
	}

	return node
}

// Bind adds an interpreter function to the node
func (n *NonTerminalNode) Bind(interpreter parsley.Interpreter) *NonTerminalNode {
	n.interpreter = interpreter
	return n
}

// Token returns with the node token
func (n *NonTerminalNode) Token() string {
	return n.token
}

// Value returns with the value of the node
func (n *NonTerminalNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	if n.interpreter == nil {
		return nil, nil
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
func (n *NonTerminalNode) ReaderPos() int {
	return n.readerPos
}

// String returns with a string representation of the node
func (n *NonTerminalNode) String() string {
	return fmt.Sprintf("NT{%s, %s, %d, %d}", n.token, n.children, n.pos, n.readerPos)
}
