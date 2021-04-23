// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// StringNode is a leaf node in the AST
type StringNode struct {
	schema    interface{}
	value     string
	pos       parsley.Pos
	readerPos parsley.Pos
}

// NewStringNode creates a new StringNode instance
func NewStringNode(schema interface{}, value string, pos parsley.Pos, readerPos parsley.Pos) *StringNode {
	return &StringNode{
		schema:    schema,
		value:     value,
		pos:       pos,
		readerPos: readerPos,
	}
}

// Token returns with the node token
func (s *StringNode) Token() string {
	return "STRING"
}

// Schema returns the schema for the node's value
func (s *StringNode) Schema() interface{} {
	return s.schema
}

// Value returns with the value of the node
func (s *StringNode) Value() interface{} {
	return s.value
}

// Pos returns the position
func (s *StringNode) Pos() parsley.Pos {
	return s.pos
}

// ReaderPos returns the position of the first character immediately after this node
func (s *StringNode) ReaderPos() parsley.Pos {
	return s.readerPos
}

// SetReaderPos changes the reader position
func (s *StringNode) SetReaderPos(fun func(parsley.Pos) parsley.Pos) {
	s.readerPos = fun(s.readerPos)
}

// String returns with a string representation of the node
func (s *StringNode) String() string {
	return fmt.Sprintf("%s{%v, %d..%d}", s.Token(), s.value, s.pos, s.readerPos)
}

// String matches a string literal enclosed in double quotes
func String(schema interface{}, allowBackquote bool) parser.Func {
	notFoundErr := parsley.NotFoundError("string literal")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		quote := '"'
		readerPos, found := tr.ReadRune(pos, quote)
		if !found {
			if allowBackquote {
				quote = '`'
				readerPos, found = tr.ReadRune(pos, quote)
			}
		}

		if !found {
			return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
		}

		// check for empty string
		readerPos, found = tr.ReadRune(readerPos, quote)
		if found {
			return NewStringNode(schema, "", pos, readerPos), data.EmptyIntSet, nil
		}

		var value []byte
		if quote == '`' {
			readerPos, value = tr.ReadRegexp(readerPos, "[^`]+")
		} else {
			readerPos, value = tr.Readf(readerPos, unquoteString)
		}

		readerPos, found = tr.ReadRune(readerPos, quote)
		if !found {
			return nil, data.EmptyIntSet, parsley.NewErrorf(readerPos, "was expecting '%s'", string(quote))
		}
		return NewStringNode(schema, string(value), pos, readerPos), data.EmptyIntSet, nil
	})
}

func unquoteString(b []byte) ([]byte, int) {
	i := 0
	for {
		if i >= len(b) {
			return b, len(b)
		}
		if b[i] == '\r' || b[i] == '\n' {
			return b[0:i], i
		}
		if b[i] == '"' {
			return b[0:i], i
		} else if b[i] == '\\' || b[i] >= utf8.RuneSelf {
			break
		}
		i++
	}

	str := string(b[i:])
	var tail string
	var res = make([]byte, 0, i)
	res = append(res, b[0:i]...)
	var err error
	var ch rune
	for {
		if str == "" {
			break
		}
		ch, _, tail, err = strconv.UnquoteChar(str, '"')
		if err != nil {
			break
		}
		res = append(res, string(ch)...)
		str = tail
	}
	return res, len(b) - len(str)
}
