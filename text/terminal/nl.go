// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"errors"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// NL matches one or more new lines
// It also reads any whitespaces before the first new line and between new lines
func NL() parsley.Parser {
	notFoundErr := errors.New("was expecting new line")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		var resPos parsley.Pos
		var resReaderPos parsley.Pos

		nlCount := 0
		for {
			pos = tr.SkipWhitespaces(pos, text.WsSpaces)
			readerPos, found := tr.ReadRune(pos, '\n')
			if !found {
				break
			}
			if resPos == 0 {
				resPos = pos
			}
			resReaderPos = readerPos
			pos = readerPos
			nlCount++
		}

		if nlCount > 0 {
			val := "\n"
			if nlCount > 1 {
				val = "\n\n"
			}
			return ast.NewTerminalNode("NL", val, resPos, resReaderPos), data.EmptyIntSet, nil
		}

		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
