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
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// Char matches a character literal enclosed in single quotes
func Char() *parser.NamedFunc {
	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, parsley.Node, parsley.Error) {
		tr := r.(*text.Reader)
		readerPos, found := tr.ReadRune(pos, '\'')
		if !found {
			return data.EmptyIntSet, nil, nil
		}

		readerPos, res := tr.ReadRegexp(
			readerPos, `\\[abfnrtv']|\\x[0-9a-fA-F]{2,2}|\\u[0-9a-fA-F]{4,4}|\\U[0-9a-fA-F]{8,8}|[^']`,
		)
		if res == nil {
			return data.EmptyIntSet, nil, parsley.NewError(readerPos, "was expecting one character")
		}

		readerPos, found = tr.ReadRune(readerPos, '\'')
		if !found {
			return data.EmptyIntSet, nil, parsley.NewError(readerPos, "was expecting \"'\"")
		}

		value, _, tail, err := strconv.UnquoteChar(string(res), '\'')
		if tail != "" || err != nil {
			return data.EmptyIntSet, nil, parsley.NewError(readerPos, "invalid character value")
		}

		return data.EmptyIntSet, ast.NewTerminalNode("CHAR", value, pos, readerPos), nil
	}).WithName("char value")
}
