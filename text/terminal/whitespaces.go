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

// Whitespaces matches one or more spaces or tabs. If newLine is true it also matches \n and \f characters.
func Whitespaces(wsMode text.WsMode) parsley.Parser {
	if wsMode == text.WsNone {
		return parser.Nil()
	}
	var name string
	if wsMode == text.WsSpaces {
		name = "spaces or tabs"
	} else {
		name = "spaces, tabs or newline"
	}
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos := tr.SkipWhitespaces(pos, wsMode); readerPos > pos {
			return ast.NilNode(readerPos), data.EmptyIntSet
		}

		return nil, data.EmptyIntSet
	}).WithName(name)
}
