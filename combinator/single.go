// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
)

// Single will change the result of p if it returns with a non terminal node
// with only one child. In this case directly the child will returned.
func Single(p parsley.Parser) parser.Func {
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		res, cp, err := p.Parse(ctx, leftRecCtx, pos)
		if err != nil {
			return nil, cp, err
		}

		if node, ok := res.(parsley.NonTerminalNode); ok {
			if len(node.Children()) == 1 {
				return node.Children()[0], cp, nil
			}
		}

		return res, cp, nil
	})
}
