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
func Memoize(p parsley.Parser) *parser.NamedFunc {
	parserIndex := int(atomic.AddInt32(&nextParserIndex, 1))
	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		if result, found := h.GetResult(parserIndex, pos, leftRecCtx); found {
			return result.CurtailingParsers, result.Nodes, result.Err
		}

		if leftRecCtx.Get(parserIndex) > r.Remaining(pos)+1 {
			return data.NewIntSet(parserIndex), nil, nil
		}

		cp, nodes, err := p.Parse(h, leftRecCtx.Inc(parserIndex), r, pos)
		leftRecCtx = leftRecCtx.Filter(cp)

		res := &parsley.Result{
			LeftRecCtx:        leftRecCtx,
			CurtailingParsers: cp,
			Nodes:             nodes,
			Err:               err,
		}
		h.SaveResult(parserIndex, pos, res)

		return cp, nodes, err
	}).WithName(p.Name())
}
