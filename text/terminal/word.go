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

// Word matches the given word
func Word(word string, value interface{}) *parser.NamedFunc {
	if word == "" {
		panic("Word() should not be called with empty word")
	}

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, found := tr.MatchWord(pos, word); found {
			return ast.NewTerminalNode("WORD", value, pos, readerPos), data.EmptyIntSet
		}
		return nil, data.EmptyIntSet
	}).WithName(fmt.Sprintf("%q", word))
}
