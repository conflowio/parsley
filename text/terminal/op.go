// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"fmt"
	"strconv"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// OpNode is a leaf node in the AST
type OpNode struct {
	value     string
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewOpNode creates a new OpNode instance
func NewOpNode(value string, pos parsley.Pos, readerPos parsley.Pos) *OpNode {
	return &OpNode{
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (o *OpNode) Token() string {
	return o.value
}

// Schema returns the schema for the node's value
func (o *OpNode) Schema() interface{} {
	return nil
}

// Value returns with the value of the node
func (o *OpNode) Value() interface{} {
	return o.value
}

// Pos returns the position
func (o *OpNode) Pos() parsley.Pos {
	return o.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (o *OpNode) ReaderPos() parsley.Pos {
	return o.readerPos
}

// SetReaderPos changes the reader position
func (o *OpNode) SetReaderPos(fun func(parsley.Pos) parsley.Pos) {
	o.readerPos = fun(o.readerPos)
}

// String returns with a string representation of the node
func (o *OpNode) String() string {
	return fmt.Sprintf("%s{%d..%d}", o.Token(), o.pos, o.readerPos)
}

// Op matches the given operator
func Op(op string) parser.Func {
	if op == "" {
		panic("Op() should not be called with empty string")
	}

	notFoundErr := parsley.NotFoundError(strconv.Quote(op))

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, found := tr.MatchString(pos, op); found {
			return NewOpNode(op, pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
