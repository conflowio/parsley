// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// FloatType contains the float type's name
const FloatType = "float64"

// FloatNode is a leaf node in the AST
type FloatNode struct {
	value     float64
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewFloatNode creates a new FloatNode instance
func NewFloatNode(value float64, pos parsley.Pos, readerPos parsley.Pos) *FloatNode {
	return &FloatNode{
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (f *FloatNode) Token() string {
	return "FLOAT"
}

// Type returns
func (f *FloatNode) Type() string {
	return FloatType
}

// Value returns with the value of the node
func (f *FloatNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
	return f.value, nil
}

// Pos returns the position
func (f *FloatNode) Pos() parsley.Pos {
	return f.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (f *FloatNode) ReaderPos() parsley.Pos {
	return f.readerPos
}

// SetReaderPos changes the reader position
func (f *FloatNode) SetReaderPos(fun func(parsley.Pos) parsley.Pos) {
	f.readerPos = fun(f.readerPos)
}

// String returns with a string representation of the node
func (f *FloatNode) String() string {
	return fmt.Sprintf("%s{%v, %d..%d}", f.Token(), f.value, f.pos, f.readerPos)
}

// Float matches a float literal
func Float() parser.Func {
	notFoundErr := errors.New("was expecting float value")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, result := tr.ReadRegexp(pos, "[-+]?[0-9]*\\.[0-9]+(?:[eE][-+]?[0-9]+)?"); result != nil {
			val, err := strconv.ParseFloat(string(result), 64)
			if err != nil {
				return nil, data.EmptyIntSet, parsley.NewErrorf(pos, "invalid float value")
			}
			return NewFloatNode(val, pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
