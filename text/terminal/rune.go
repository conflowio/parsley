// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"regexp"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
)

// Rune matches the given character
func Rune(char rune, token string) parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		tr := r.(*text.Reader)
		if _, pos, ok := tr.ReadMatch(regexp.QuoteMeta(string(char)), false); ok {
			return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(token, pos, char), r).AsSet(), nil
		}
		return parser.NoCurtailingParsers(), nil, parser.NewError(r.Cursor(), "was expecting '%s'", string(char))
	})
}
