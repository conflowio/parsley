// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
)

// Any tries all the given parsers independently and merges the results
func Any(parsers ...parsley.Parser) parser.Func {
	if parsers == nil {
		panic("no parsers were given")
	}

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		cp := data.EmptyIntSet
		var res parsley.Node
		var err parsley.Error
		for _, p := range parsers {
			ctx.RegisterCall()
			res2, cp2, err2 := p.Parse(ctx, leftRecCtx, pos)
			cp = cp.Union(cp2)
			res = ast.AppendNode(res, res2)
			if err2 != nil && (err == nil || err2.Pos() >= err.Pos()) {
				if err2.Pos() > pos || !parsley.IsNotFoundError(err2) {
					err = err2
				}
			}
		}

		if res == nil {
			return nil, cp, err
		}

		ctx.SetError(err)

		return res, cp, nil
	})
}
