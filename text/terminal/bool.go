// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"fmt"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// BoolType contains the boolean type's name
const BoolType = "bool"

// BoolNode is a leaf node in the AST
type BoolNode struct {
	value     bool
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewBoolNode creates a new BoolNode instance
func NewBoolNode(value bool, pos parsley.Pos, readerPos parsley.Pos) *BoolNode {
	return &BoolNode{
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (b *BoolNode) Token() string {
	return "BOOL"
}

// Type returns
func (b *BoolNode) Type() string {
	return BoolType
}

// Value returns with the value of the node
func (b *BoolNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
	return b.value, nil
}

// Pos returns the position
func (b *BoolNode) Pos() parsley.Pos {
	return b.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (b *BoolNode) ReaderPos() parsley.Pos {
	return b.readerPos
}

// SetReaderPos changes the reader position
func (b *BoolNode) SetReaderPos(fun func(parsley.Pos) parsley.Pos) {
	b.readerPos = fun(b.readerPos)
}

// String returns with a string representation of the node
func (b *BoolNode) String() string {
	return fmt.Sprintf("%s{%v, %d..%d}", b.Token(), b.value, b.pos, b.readerPos)
}

// Bool matches a bool literal: true or false
func Bool(trueStr string, falseStr string) parser.Func {
	if trueStr == "" || falseStr == "" {
		panic("Bool() should not be called with an empty true/false string")
	}

	notFoundErr := parsley.NotFoundError("boolean")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, found := tr.MatchWord(pos, trueStr); found {
			return NewBoolNode(true, pos, readerPos), data.EmptyIntSet, nil
		}
		if readerPos, found := tr.MatchWord(pos, falseStr); found {
			return NewBoolNode(false, pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
