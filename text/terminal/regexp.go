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
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// Regexp matches the given regular expression
// The name variable is used for error messages, so it should be descriptive and make sense in the sentence "was expecting %s".
// The includeWhitespaces variable should be true if the reader is by default ignoring the whitespaces but you need to match those as well.
// If you are using capturing groups you can select which group to use as a value with the groupIdex variable.
func Regexp(token string, name string, regexp string, groupIndex int, wsMode text.WsMode) *parser.NamedFunc {
	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, parsley.Node, parsley.Error) {
		tr := r.(*text.Reader)
		if groupIndex == 0 {
			if readerPos, match := tr.ReadRegexp(pos, regexp); match != nil {
				readerPos = tr.SkipWhitespaces(readerPos, wsMode)
				return data.EmptyIntSet, ast.NewTerminalNode(token, string(match), r.Pos(pos), readerPos), nil
			}
		} else {
			if readerPos, matches := tr.ReadRegexpSubmatch(pos, regexp); matches != nil {
				if groupIndex >= len(matches) {
					panic(fmt.Sprintf("Capturing group %d is invalid for %s", groupIndex, regexp))
				}
				readerPos = tr.SkipWhitespaces(readerPos, wsMode)
				return data.EmptyIntSet, ast.NewTerminalNode(token, string(matches[groupIndex]), r.Pos(pos), readerPos), nil
			}
		}
		return data.EmptyIntSet, nil, nil
	}).WithName(name)
}
