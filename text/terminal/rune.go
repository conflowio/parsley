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

// Rune matches the given character
func Rune(ch rune) parsley.ParserFunc {
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		tr := r.(*text.Reader)
		if readerPos, found := tr.ReadRune(pos, ch); found {
			return data.EmptyIntSet(), []parsley.Node{ast.NewTerminalNode(string(ch), ch, r.Pos(pos), readerPos)}, nil
		}
		return data.EmptyIntSet(), nil, parsley.NewError(r.Pos(pos), "was expecting %q", string(ch))
	})
}
