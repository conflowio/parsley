// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast

import (
	"fmt"

	"github.com/opsidian/parsley/reader"
)

// EOF is the end of file token
const EOF = "EOF"

// Node represents an AST node
type Node interface {
	Token() string
	Value(ctx interface{}) (interface{}, reader.Error)
	Pos() reader.Position
}

// TerminalNode is a leaf node in the AST
type TerminalNode struct {
	token string
	pos   reader.Position
	value interface{}
}

// NewTerminalNode creates a new TerminalNode instance
func NewTerminalNode(token string, pos reader.Position, value interface{}) TerminalNode {
	return TerminalNode{
		token: token,
		pos:   pos,
		value: value,
	}
}

// Token returns with the node token
func (t TerminalNode) Token() string {
	return t.token
}

// Value returns with the value of the node
func (t TerminalNode) Value(ctx interface{}) (interface{}, reader.Error) {
	return t.value, nil
}

// Pos returns the position
func (t TerminalNode) Pos() reader.Position {
	return t.pos
}

// String returns with a string representation of the node
func (t TerminalNode) String() string {
	if t.value == nil {
		return fmt.Sprintf("T{%s, %s}", t.token, t.pos)
	}
	return fmt.Sprintf("T{%v, %s}", t.value, t.pos)
}

// NonTerminalNode represents a branch node in the AST
type NonTerminalNode struct {
	token       string
	pos         reader.Position
	children    []Node
	interpreter Interpreter
}

// NewNonTerminalNode creates a new NonTerminalNode instance
func NewNonTerminalNode(token string, children []Node, interpreter Interpreter) NonTerminalNode {
	var pos reader.Position
	for _, child := range children {
		if child != nil {
			pos = child.Pos()
			break
		}
	}
	return NonTerminalNode{
		token:       token,
		pos:         pos,
		children:    children,
		interpreter: interpreter,
	}
}

// Token returns with the node token
func (n NonTerminalNode) Token() string {
	return n.token
}

// Value returns with the value of the node
func (n NonTerminalNode) Value(ctx interface{}) (interface{}, reader.Error) {
	if n.interpreter == nil {
		return nil, nil
	}
	return n.interpreter.Eval(ctx, n.children)
}

// Pos returns the position
func (n NonTerminalNode) Pos() reader.Position {
	return n.pos
}

// Children returns with the children
func (n NonTerminalNode) Children() []Node {
	return n.children
}

// String returns with a string representation of the node
func (n NonTerminalNode) String() string {
	return fmt.Sprintf("NT{%s, %s}", n.token, n.children)
}
