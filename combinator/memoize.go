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
	"github.com/opsidian/parsley/reader"
)

var nextParserIndex int32

// Memoize handles result cache and curtailing left recursion
func Memoize(p parser.Parser) parser.Func {
	parserIndex := int(atomic.AddInt32(&nextParserIndex, 1))
	return parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		cp, rs, err, found := h.GetResults(parserIndex, r.Cursor().Pos(), leftRecCtx)
		if found {
			return cp, rs, err
		}

		if leftRecCtx.Get(parserIndex) > r.Remaining()+1 {
			return data.NewIntSet(parserIndex), nil, nil
		}

		cp, rs, err = p.Parse(h, leftRecCtx.Inc(parserIndex), r)
		leftRecCtx = leftRecCtx.Filter(cp)

		h.RegisterResults(parserIndex, r.Cursor().Pos(), cp, rs, err, leftRecCtx)

		return cp, rs, err
	})
}
