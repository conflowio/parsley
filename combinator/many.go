// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Many applies the  parser zero or more times
func Many(nodeBuilder ast.NodeBuilder, p parser.Parser) parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return newRecursive(nodeBuilder, func(i int) parser.Parser { return p }, 0, -1).Parse(leftRecCtx, r)
	})
}

// Many1 applies the parser one or more times
func Many1(nodeBuilder ast.NodeBuilder, p parser.Parser) parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return newRecursive(nodeBuilder, func(i int) parser.Parser { return p }, 1, -1).Parse(leftRecCtx, r)
	})
}
