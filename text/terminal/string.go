// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"strconv"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// String matches a string literal enclosed in double quotes
func String(allowBackquote bool) parsley.ParserFunc {
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		tr := r.(*text.Reader)
		quote := '"'
		readerPos, found := tr.ReadRune(pos, quote)
		if !found {
			if allowBackquote {
				quote = '`'
				readerPos, found = tr.ReadRune(pos, quote)
			}
		}

		if !found {
			return data.EmptyIntSet, nil, parsley.NewError(r.Pos(pos), "was expecting string literal")
		}

		// check for empty string
		readerPos, found = tr.ReadRune(readerPos, quote)
		if found {
			return data.EmptyIntSet, []parsley.Node{ast.NewTerminalNode("STRING", "", r.Pos(pos), readerPos)}, nil
		}

		var value []byte
		if quote == '`' {
			readerPos, value = tr.ReadRegexp(readerPos, "[^`]+")
		} else {
			readerPos, value = tr.Readf(readerPos, unquoteString)
		}

		readerPos, found = tr.ReadRune(readerPos, quote)
		if !found {
			return data.EmptyIntSet, nil, parsley.NewError(r.Pos(readerPos), "was expecting '%s'", string(quote))
		}
		return data.EmptyIntSet, []parsley.Node{ast.NewTerminalNode("STRING", string(value), r.Pos(pos), readerPos)}, nil
	})
}

func unquoteString(b []byte) ([]byte, int) {
	str := string(b)
	var tail, res string
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
		res += string(ch)
		str = tail
	}
	return []byte(res), len(b) - len(str)
}
