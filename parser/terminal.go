// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/reader"
)

// Empty always matches and returns with an empty node result
// When using Empty you should not forget to handle for nil nodes in your node builders and/or interpreters.
func Empty() Func {
	return Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, ResultSet, Error) {
		return NoCurtailingParsers(), NewResult(nil, r).AsSet(), nil
	})
}

// End matches the end of the input
func End() Func {
	return Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, ResultSet, Error) {
		if r.IsEOF() {
			return NoCurtailingParsers(), NewResult(ast.NewTerminalNode(ast.EOF, r.Cursor(), nil), r).AsSet(), nil
		}
		return NoCurtailingParsers(), nil, NewError(r.Cursor(), "was expecting the end of input")
	})
}
