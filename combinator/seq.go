// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

// Seq tries to apply all parsers after each other matching effectively a sequence of tokens and returns with all
// combinations of the results. Only matches are returned where all parsers were applied successfully.
func Seq(nodeBuilder parsley.NodeBuilder, parsers ...parsley.Parser) parsley.ParserFunc {
	if parsers == nil {
		panic("No parsers were given")
	}
	lookup := func(i int) parsley.Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		l := len(parsers)
		return newRecursive(nodeBuilder, lookup, l, l).Parse(h, leftRecCtx, r, pos)
	})
}

// SeqTry tries to apply all parsers after each other matching effectively the longest possible sequences of
// tokens and returns with all combinations of the results.
func SeqTry(nodeBuilder parsley.NodeBuilder, min int, parsers ...parsley.Parser) parsley.ParserFunc {
	if parsers == nil {
		panic("No parsers were given")
	}
	lookup := func(i int) parsley.Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		return newRecursive(nodeBuilder, lookup, min, len(parsers)).Parse(h, leftRecCtx, r, pos)
	})
}
