// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// Any tries all the given parsers independently and merges the results
func Any(name string, parsers ...parsley.Parser) *parser.NamedFunc {
	if parsers == nil {
		panic("no parsers were given")
	}

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		cp := data.EmptyIntSet
		var res parsley.Node
		for _, p := range parsers {
			ctx.RegisterCall()
			res2, cp2 := p.Parse(ctx, leftRecCtx, pos)
			cp = cp.Union(cp2)
			res = ast.AppendNode(res, res2)
		}
		return res, cp
	}).WithName(name)
}
