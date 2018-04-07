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
	"github.com/opsidian/parsley/text/terminal"
)

// Trim reads optional whitespaces and returns only the result of p
func Trim(p parsley.Parser, newLine bool) *parser.NamedFunc {
	ws := terminal.Whitespaces(newLine)
	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, parsley.Node, parsley.Error) {
		_, node, _ := ws.Parse(h, leftRecCtx, r, pos)
		if node != nil {
			pos = node.ReaderPos()
		}
		return p.Parse(h, leftRecCtx, r, pos)
	}).WithName(p.Name())
}
