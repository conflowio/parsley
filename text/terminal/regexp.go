// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
)

// Regexp matches the given regular expression
// The desc variable is used for error messages, so it should be descriptive and make sense in the sentence "was expecting %s".
// The includeWhitespaces variable should be true if the reader is by default ignoring the whitespaces but you need to match those as well.
// If you are using capturing groups you can select which group to use as a value with the groupIdex variable.
func Regexp(desc string, regexp string, includeWhitespaces bool, groupIndex int, token string) parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		tr := r.(*text.Reader)
		if matches, pos, ok := tr.ReadMatch(regexp, includeWhitespaces); ok {
			if groupIndex >= len(matches) {
				panic(fmt.Sprintf("Capturing group %d is invalid for %s", groupIndex, regexp))
			}
			if token == "" {
				token = matches[groupIndex]
			}
			return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(token, pos, matches[groupIndex]), r).AsSet(), nil
		}
		return parser.NoCurtailingParsers(), nil, reader.NewError(r.Cursor(), "was expecting %s", desc)
	})
}
