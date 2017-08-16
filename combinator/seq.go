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

// Seq tries to apply all parsers after each other matching effectively a sequence of tokens and returns with all
// combinations of the results. Only matches are returned where all parsers were applied successfully.
func Seq(nodeBuilder ast.NodeBuilder, parsers ...parser.Parser) parser.Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	lookup := func(i int) parser.Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		l := len(parsers)
		return newRecursive(nodeBuilder, lookup, l, l).Parse(leftRecCtx, r)
	})
}

// SeqTry tries to apply all parsers after each other matching effectively the longest possible sequences of
// tokens and returns with all combinations of the results.
func SeqTry(nodeBuilder ast.NodeBuilder, min int, parsers ...parser.Parser) parser.Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	lookup := func(i int) parser.Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return newRecursive(nodeBuilder, lookup, min, len(parsers)).Parse(leftRecCtx, r)
	})
}
