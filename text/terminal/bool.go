// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// Bool matches a bool literal: true or false
func Bool(trueStr string, falseStr string) parsley.ParserFunc {
	if trueStr == "" || falseStr == "" {
		panic("Bool() should not be called with an empty true/false string")
	}

	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		tr := r.(*text.Reader)
		if readerPos, found := tr.MatchWord(pos, trueStr); found {
			return data.EmptyIntSet(), []parsley.Node{ast.NewTerminalNode("BOOL", true, r.Pos(pos), readerPos)}, nil
		}
		if readerPos, found := tr.MatchWord(pos, falseStr); found {
			return data.EmptyIntSet(), []parsley.Node{ast.NewTerminalNode("BOOL", false, r.Pos(pos), readerPos)}, nil
		}
		return data.EmptyIntSet(), nil, parsley.NewError(r.Pos(pos), "was expecting %s or %s", trueStr, falseStr)
	})
}
