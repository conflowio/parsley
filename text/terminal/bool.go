// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// Bool matches a bool literal: true or false
func Bool(trueStr string, falseStr string) *parser.NamedFunc {
	if trueStr == "" || falseStr == "" {
		panic("Bool() should not be called with an empty true/false string")
	}

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, found := tr.MatchWord(pos, trueStr); found {
			return ast.NewTerminalNode("BOOL", true, pos, readerPos), data.EmptyIntSet
		}
		if readerPos, found := tr.MatchWord(pos, falseStr); found {
			return ast.NewTerminalNode("BOOL", false, pos, readerPos), data.EmptyIntSet
		}
		return nil, data.EmptyIntSet
	}).WithName(fmt.Sprintf("%s or %s", trueStr, falseStr))
}
