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
func Choice(parsers ...parsley.Parser) parser.Func {
	if parsers == nil {
		panic("No parsers were given")
	}

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		cp := data.EmptyIntSet
		var err parsley.Error
		for _, p := range parsers {
			ctx.RegisterCall()
			node, cp2, err2 := p.Parse(ctx, leftRecCtx, pos)
			cp = cp.Union(cp2)

			if err2 != nil && (err == nil || err2.Pos() >= err.Pos()) {
				if err2.Pos() > pos || !parsley.IsNotFoundError(err2) {
					err = err2
				}
			}
			if node != nil {
				ctx.SetError(err)
				return node, cp, nil
			}
		}

		return nil, cp, err
	})
}
