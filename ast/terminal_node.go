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

// TerminalNode is a leaf node in the AST
type TerminalNode struct {
	token     string
	valueType string
	value     interface{}
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewTerminalNode creates a new TerminalNode instance
func NewTerminalNode(
	token string,
	value interface{},
	valueType string,
	pos parsley.Pos,
	readerPos parsley.Pos,
) *TerminalNode {
	return &TerminalNode{
		token:     token,
		value:     value,
		valueType: valueType,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (t *TerminalNode) Token() string {
	return t.token
}

// Type returns with the type of the node's value
func (t *TerminalNode) Type() string {
	return t.valueType
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
	if t.valueType == "" {
		return fmt.Sprintf("%s{%v, %d..%d}", t.token, t.value, t.pos, t.readerPos)
	}

	return fmt.Sprintf("%s{<%s> %v, %d..%d}", t.token, t.valueType, t.value, t.pos, t.readerPos)
}
