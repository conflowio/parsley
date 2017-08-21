// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"fmt"
	"strconv"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
)

// Char matches a character literal enclosed in single quotes
func Char() parser.Func {
	return parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		tr := r.(*text.Reader)
		if matches, pos, ok := tr.ReadMatch("'(.|\\\\[abfnrtv]|\\\\x[0-9a-fA-F]{2,2}|\\\\u[0-9a-fA-F]{4,4}|\\\\U[0-9a-fA-F]{8,8})'", false); ok {
			match := matches[1]
			value, _, tail, err := strconv.UnquoteChar(match, '"')
			if tail != "" {
				// This should never happen
				panic(fmt.Sprintf("Unprocessed string segment: %s", tail))
			}
			if err == nil {
				return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("CHAR", pos, value), r).AsSet(), nil
			}
		}
		return parser.NoCurtailingParsers(), nil, reader.NewError(r.Cursor(), "was expecting char literal")
	})
}
