// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// Choice tries to apply the given parsers until one of them succeeds
func Choice(name string, parsers ...parsley.Parser) *parser.NamedFunc {
	if parsers == nil {
		panic("No parsers were given")
	}

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		cp := data.EmptyIntSet
		for _, p := range parsers {
			ctx.RegisterCall()
			node, cp2 := p.Parse(ctx, leftRecCtx, pos)
			cp = cp.Union(cp2)
			if node != nil {
				return node, cp
			}
		}
		return nil, cp
	}).WithName(name)
}
