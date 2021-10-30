// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"strconv"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
)

// Rune matches the given character
func Rune(ch rune) parser.Func {
	notFoundErr := parsley.NotFoundError(strconv.Quote(string(ch)))

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, found := tr.ReadRune(pos, ch); found {
			return ast.NewTerminalNode(nil, string(ch), ch, pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
