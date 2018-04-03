// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

// Many applies the  parser zero or more times
func Many(nodeBuilder parsley.NodeBuilder, p parsley.Parser) parsley.ParserFunc {
	f := func(i int) parsley.Parser { return p }
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		return newRecursive(nodeBuilder, f, 0, -1).Parse(h, leftRecCtx, r, pos)
	})
}

// Many1 applies the parser one or more times
func Many1(nodeBuilder parsley.NodeBuilder, p parsley.Parser) parsley.ParserFunc {
	f := func(i int) parsley.Parser { return p }
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		return newRecursive(nodeBuilder, f, 1, -1).Parse(h, leftRecCtx, r, pos)
	})
}
