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

// IntegerType contains the integer type's name
const IntegerType = "int64"

// IntegerNode is a leaf node in the AST
type IntegerNode struct {
	value     int64
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewIntegerNode creates a new IntegerNode instance
func NewIntegerNode(value int64, pos parsley.Pos, readerPos parsley.Pos) *IntegerNode {
	return &IntegerNode{
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (i *IntegerNode) Token() string {
	return "INTEGER"
}

// Type returns
func (i *IntegerNode) Type() string {
	return IntegerType
}

// Value returns with the value of the node
func (i *IntegerNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
	return i.value, nil
}

// Pos returns the position
func (i *IntegerNode) Pos() parsley.Pos {
	return i.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (i *IntegerNode) ReaderPos() parsley.Pos {
	return i.readerPos
}

// SetReaderPos changes the reader position
func (i *IntegerNode) SetReaderPos(fun func(parsley.Pos) parsley.Pos) {
	i.readerPos = fun(i.readerPos)
}

// String returns with a string representation of the node
func (i *IntegerNode) String() string {
	return fmt.Sprintf("%s{%v, %d..%d}", i.Token(), i.value, i.pos, i.readerPos)
}

// Integer matches all integer numbers and zero with an optional -/+ sign
func Integer() parser.Func {
	notFoundErr := parsley.NewNotFoundError("integer value")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, result := tr.ReadRegexp(pos, "[-+]?(?:[1-9][0-9]*|0[xX][0-9a-fA-F]+|0[0-7]*)"); result != nil {
			if _, isFloat := tr.ReadRune(readerPos, '.'); isFloat {
				return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
			}
			intValue, err := strconv.ParseInt(string(result), 0, 64)
			if err != nil {
				// This should never happen
				panic(fmt.Sprintf("Could not convert %s to integer", string(result)))
			}
			return NewIntegerNode(intValue, pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
