// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// Nil matches a nil literal
func Nil(nilStr string) *parser.NamedFunc {
	if nilStr == "" {
		panic("Nil() should not be called with an empty nil string")
	}

	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, parsley.Node, parsley.Error) {
		tr := r.(*text.Reader)
		if readerPos, found := tr.MatchWord(pos, nilStr); found {
			return data.EmptyIntSet, ast.NewTerminalNode("NIL", nil, pos, readerPos), nil
		}

		return data.EmptyIntSet, nil, nil
	}).WithName(nilStr)
}
