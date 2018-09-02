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
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		tr := ctx.Reader().(*text.Reader)
		readerPos, found := tr.ReadRune(pos, '\'')
		if !found {
			return nil, data.EmptyIntSet
		}

		readerPos, res := tr.ReadRegexp(
			readerPos, `\\[abfnrtv']|\\x[0-9a-fA-F]{2,2}|\\u[0-9a-fA-F]{4,4}|\\U[0-9a-fA-F]{8,8}|[^']`,
		)
		if res == nil {
			ctx.SetErrorf(readerPos, "was expecting one character")
			return nil, data.EmptyIntSet
		}

		readerPos, found = tr.ReadRune(readerPos, '\'')
		if !found {
			ctx.SetErrorf(readerPos, "was expecting \"'\"")
			return nil, data.EmptyIntSet
		}

		value, _, tail, err := strconv.UnquoteChar(string(res), '\'')
		if tail != "" || err != nil {
			ctx.SetErrorf(readerPos, "invalid character value")
			return nil, data.EmptyIntSet
		}

		return ast.NewTerminalNode("CHAR", value, pos, readerPos), data.EmptyIntSet
	}).WithName("char value")
}
