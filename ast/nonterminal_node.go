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

// NonTerminalNode represents a branch node in the AST
type NonTerminalNode struct {
	token       string
	valueType   string
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

// Type returns with the type of the node's value
func (n *NonTerminalNode) Type() string {
	return n.valueType
}

// Value returns with the value of the node
func (n *NonTerminalNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
	if n.interpreter == nil {
		panic("missing interpreter for node")
	}

	return n.interpreter.Eval(userCtx, n)
}

// Transform runs the given transformer on all children and returns the original node
func (n *NonTerminalNode) Transform(userCtx interface{}) (parsley.Node, parsley.Error) {
	if n.interpreter != nil {
		switch i := n.interpreter.(type) {
		case parsley.NodeTransformer:
			return i.TransformNode(userCtx, n)
		}
	}

	var err parsley.Error
	for i, child := range n.children {
		if n.children[i], err = parsley.Transform(userCtx, child); err != nil {
			return nil, err
		}
	}

	return n, nil
}

// StaticCheck runs a static analysis if the interpreter has static analysis capabilities
func (n *NonTerminalNode) StaticCheck(userCtx interface{}) parsley.Error {
	if n.interpreter != nil {
		switch i := n.interpreter.(type) {
		case parsley.StaticChecker:
			valueType, err := i.StaticCheck(userCtx, n)
			if err != nil {
				return err
			}
			n.valueType = valueType
		}
	}

	return nil
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
	if n.valueType == "" {
		return fmt.Sprintf("%s{%s, %d..%d}", n.token, n.children, n.pos, n.readerPos)
	}
	return fmt.Sprintf("%s{<%s> %s, %d..%d}", n.token, n.valueType, n.children, n.pos, n.readerPos)
}
