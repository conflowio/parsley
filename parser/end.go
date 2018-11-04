// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"errors"
	"fmt"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

// EOF is the end-of-file token
const EOF = "EOF"

// EndNode represents an eof node
type EndNode parsley.Pos

// Token returns with EOF
func (e EndNode) Token() string {
	return EOF
}

// Type returns with an empty type
func (e EndNode) Type() string {
	return ""
}

// Value returns with nil
func (e EndNode) Value(ctx interface{}) (interface{}, parsley.Error) {
	return nil, nil
}

// Pos returns with the position
func (e EndNode) Pos() parsley.Pos {
	return parsley.Pos(e)
}

// ReaderPos returns the reader position
func (e EndNode) ReaderPos() parsley.Pos {
	return parsley.Pos(e)
}

// SetReaderPos changes the reader position
func (e EndNode) SetReaderPos(func(parsley.Pos) parsley.Pos) {
}

// String returns with a string representation of the node
func (e EndNode) String() string {
	return fmt.Sprintf("%s{%d}", e.Token(), e.Pos())
}

// End matches the end of the input
func End() Func {
	notFoundErr := errors.New("was expecting the end of input")

	return Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		if ctx.Reader().IsEOF(pos) {
			return EndNode(pos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
