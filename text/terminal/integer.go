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
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// Integer matches all integer numbers and zero with an optional -/+ sign
func Integer() *parser.NamedFunc {
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, result := tr.ReadRegexp(pos, "[-+]?(?:[1-9][0-9]*|0[xX][0-9a-fA-F]+|0[0-7]*)"); result != nil {
			if _, isFloat := tr.ReadRune(readerPos, '.'); isFloat {
				return nil, data.EmptyIntSet
			}
			intValue, err := strconv.ParseInt(string(result), 0, 64)
			if err != nil {
				// This should never happen
				panic(fmt.Sprintf("Could not convert %s to integer", string(result)))
			}
			return ast.NewTerminalNode("INT", intValue, pos, readerPos), data.EmptyIntSet
		}
		return nil, data.EmptyIntSet
	}).WithName("integer value")
}
