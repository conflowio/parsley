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

// CharType contains the char type's name
const CharType = "rune"

// CharNode is a leaf node in the AST
type CharNode struct {
	value     rune
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewCharNode creates a new CharNode instance
func NewCharNode(value rune, pos parsley.Pos, readerPos parsley.Pos) *CharNode {
	return &CharNode{
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (c *CharNode) Token() string {
	return "CHAR"
}

// Type returns
func (c *CharNode) Type() string {
	return CharType
}

// Value returns with the value of the node
func (c *CharNode) Value(userCtx interface{}) (interface{}, parsley.Error) {
	return c.value, nil
}

// Pos returns the position
func (c *CharNode) Pos() parsley.Pos {
	return c.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (c *CharNode) ReaderPos() parsley.Pos {
	return c.readerPos
}

// SetReaderPos changes the reader position
func (c *CharNode) SetReaderPos(fun func(parsley.Pos) parsley.Pos) {
	c.readerPos = fun(c.readerPos)
}

// String returns with a string representation of the node
func (c *CharNode) String() string {
	return fmt.Sprintf("%s{%v, %d..%d}", c.Token(), c.value, c.pos, c.readerPos)
}

// Char matches a character literal enclosed in single quotes
func Char() parser.Func {
	notFoundErr := parsley.NewNotFoundError("char literal")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		readerPos, found := tr.ReadRune(pos, '\'')
		if !found {
			return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
		}

		readerPos, res := tr.ReadRegexp(
			readerPos, `\\[abfnrtv']|\\x[0-9a-fA-F]{2,2}|\\u[0-9a-fA-F]{4,4}|\\U[0-9a-fA-F]{8,8}|[^']`,
		)
		if res == nil {
			return nil, data.EmptyIntSet, parsley.NewErrorf(readerPos, "was expecting one character")
		}

		readerPos, found = tr.ReadRune(readerPos, '\'')
		if !found {
			return nil, data.EmptyIntSet, parsley.NewErrorf(readerPos, "was expecting \"'\"")
		}

		value, _, tail, err := strconv.UnquoteChar(string(res), '\'')
		if tail != "" || err != nil {
			return nil, data.EmptyIntSet, parsley.NewErrorf(readerPos, "invalid character value")
		}

		return NewCharNode(value, pos, readerPos), data.EmptyIntSet, nil
	})
}
