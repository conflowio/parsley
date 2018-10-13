// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"errors"
	"strconv"
	"unicode/utf8"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// String matches a string literal enclosed in double quotes
func String(allowBackquote bool) parser.Func {
	notFoundErr := errors.New("was expecting string literal")

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
			return ast.NewTerminalNode("STRING", "", pos, readerPos), data.EmptyIntSet, nil
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
		return ast.NewTerminalNode("STRING", string(value), pos, readerPos), data.EmptyIntSet, nil
	})
}

func unquoteString(b []byte) ([]byte, int) {
	i := 0
	for {
		if i >= len(b) {
			return b, len(b)
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
