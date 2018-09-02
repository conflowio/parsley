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

// Float matches a float literal
func Float() *parser.NamedFunc {
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, result := tr.ReadRegexp(pos, "[-+]?[0-9]*\\.[0-9]+(?:[eE][-+]?[0-9]+)?"); result != nil {
			val, err := strconv.ParseFloat(string(result), 64)
			if err != nil {
				ctx.SetErrorf(pos, "invalid float value encountered")
				return nil, data.EmptyIntSet
			}
			return ast.NewTerminalNode("FLOAT", val, pos, readerPos), data.EmptyIntSet
		}
		return nil, data.EmptyIntSet
	}).WithName("float value")
}
