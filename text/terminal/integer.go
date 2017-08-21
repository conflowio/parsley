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

// Integer matches all integer numbers and zero with an optional -/+ sign
func Integer() parser.Func {
	return parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		tr := r.(*text.Reader)
		if matches, pos, ok := tr.ReadMatch("[-+]?(?:[1-9][0-9]*|0[xX][0-9a-fA-F]+|0[0-7]*)", false); ok {
			intValue, err := strconv.ParseInt(matches[0], 0, 0)
			if err != nil {
				// This should never happen
				panic(fmt.Sprintf("Could not convert %s to integer", matches[0]))
			}
			return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("INT", pos, int(intValue)), r).AsSet(), nil
		}
		return parser.NoCurtailingParsers(), nil, reader.NewError(r.Cursor(), "was expecting integer value")
	})
}
