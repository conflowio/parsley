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

// Whitespaces matches one or more spaces or tabs. If newLine is true it also matches \n and \f characters.
func Whitespaces(wsMode text.WsMode) parsley.Parser {
	if wsMode == text.WsNone {
		return parser.Nil()
	}

	var notFoundErr error
	if wsMode == text.WsSpacesForceNl {
		notFoundErr = errors.New("was expecting a new line")
	} else if wsMode == text.WsSpaces {
		notFoundErr = errors.New("was expecting spaces or tabs")
	} else {
		notFoundErr = errors.New("was expecting spaces, tabs or a newline character")
	}

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, ok := tr.SkipWhitespaces(pos, wsMode); readerPos > pos && ok {
			return ast.NilNode(readerPos), data.EmptyIntSet, nil
		}

		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
