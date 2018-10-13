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

// Nil matches a nil literal
func Nil(nilStr string) parser.Func {
	if nilStr == "" {
		panic("Nil() should not be called with an empty nil string")
	}

	notFoundErr := fmt.Errorf("was expecting %s", nilStr)

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, found := tr.MatchWord(pos, nilStr); found {
			return ast.NewTerminalNode("NIL", nil, pos, readerPos), data.EmptyIntSet, nil
		}

		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
