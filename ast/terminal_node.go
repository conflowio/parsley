// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast

import (
	"fmt"

	"github.com/conflowio/parsley/parsley"
)

// TerminalNode is a leaf node in the AST
type TerminalNode struct {
	schema    interface{}
	token     string
	value     interface{}
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewTerminalNode creates a new TerminalNode instance
func NewTerminalNode(
	schema interface{},
	token string,
	value interface{},
	pos parsley.Pos,
	readerPos parsley.Pos,
) *TerminalNode {
	return &TerminalNode{
		schema:    schema,
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

// Schema returns the schema for the node's value
func (t *TerminalNode) Schema() interface{} {
	return t.schema
}

// Value returns with the value of the node
func (t *TerminalNode) Value() interface{} {
	return t.value
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
