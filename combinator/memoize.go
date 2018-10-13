// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"sync/atomic"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

var nextParserIndex int32

// Memoize handles result cache and curtailing left recursion
func Memoize(p parsley.Parser) parser.Func {
	parserIndex := int(atomic.AddInt32(&nextParserIndex, 1))
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		if result, found := ctx.ResultCache().Get(parserIndex, pos, leftRecCtx); found {
			return result.Node, result.CurtailingParsers, result.Error
		}

		if leftRecCtx.Get(parserIndex) > ctx.Reader().Remaining(pos)+1 {
			return nil, data.NewIntSet(parserIndex), nil
		}

		node, cp, err := p.Parse(ctx, leftRecCtx.Inc(parserIndex), pos)
		leftRecCtx = leftRecCtx.Filter(cp)

		res := &parsley.Result{
			LeftRecCtx:        leftRecCtx,
			CurtailingParsers: cp,
			Error:             err,
			Node:              node,
		}
		ctx.ResultCache().Save(parserIndex, pos, res)

		return node, cp, err
	})
}
