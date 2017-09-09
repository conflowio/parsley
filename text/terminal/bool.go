// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
)

// Bool matches a bool literal: true or false
func Bool() parser.Func {
	return parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		tr := r.(*text.Reader)
		if matches, pos, ok := tr.ReadMatch("true|false", false); ok {
			if _, ok := tr.PeekMatch("\\w+"); ok {
				return parser.NoCurtailingParsers(), nil, reader.NewError(pos, "was expecting boolean")
			}
			val := false
			if matches[0] == "true" {
				val = true
			}
			return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("BOOL", pos, val), r).AsSet(), nil
		}
		return parser.NoCurtailingParsers(), nil, reader.NewError(r.Cursor(), "was expecting boolean")
	})
}
