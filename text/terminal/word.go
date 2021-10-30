// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"strconv"
	"strings"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
)

// Word matches the given word
func Word(schema interface{}, word string, value interface{}) parser.Func {
	if word == "" {
		panic("Word() should not be called with empty word")
	}

	notFoundErr := parsley.NotFoundError(strconv.Quote(word))
	token := strings.ToUpper(word)

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, found := tr.MatchWord(pos, word); found {
			return ast.NewTerminalNode(schema, token, value, pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
