// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

// Empty always matches and returns with an empty node result
// When using Empty you should not forget to handle for nil nodes in your node builders and/or interpreters.
func Empty() parsley.ParserFunc {
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		return data.EmptyIntSet, []parsley.Node{nil}, nil
	})
}

// End matches the end of the input
func End() parsley.ParserFunc {
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		if r.IsEOF(pos) {
			return data.EmptyIntSet, []parsley.Node{ast.NewTerminalNode(ast.EOF, nil, r.Pos(pos), pos)}, nil
		}
		return data.EmptyIntSet, nil, parsley.NewError(r.Pos(pos), "was expecting the end of input")
	})
}
