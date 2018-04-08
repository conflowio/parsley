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
		return parser.Empty()
	}

	var name string
	if wsMode == text.WsSpaces {
		name = "spaces or tabs"
	} else {
		name = "spaces, tabs or newline"
	}
	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, parsley.Node, parsley.Error) {
		tr := r.(*text.Reader)
		if readerPos := tr.SkipWhitespaces(pos, wsMode); readerPos > pos {
			return data.EmptyIntSet, ast.EmptyNode(readerPos), nil
		}

		return data.EmptyIntSet, nil, nil
	}).WithName(name)
}
