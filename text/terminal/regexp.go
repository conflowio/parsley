// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"fmt"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
)

// Regexp matches the given regular expression
// The name variable is used for error messages, so it should be descriptive and make sense in the sentence "was expecting %s".
// The includeWhitespaces variable should be true if the reader is by default ignoring the whitespaces but you need to match those as well.
// If you are using capturing groups you can select which group to use as a value with the groupIdex variable.
func Regexp(schema interface{}, token string, name string, regexp string, groupIndex int) parser.Func {
	notFoundErr := parsley.NotFoundError(name)

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if groupIndex == 0 {
			if readerPos, match := tr.ReadRegexp(pos, regexp); match != nil {
				return ast.NewTerminalNode(schema, token, string(match), pos, readerPos), data.EmptyIntSet, nil
			}
		} else {
			if readerPos, matches := tr.ReadRegexpSubmatch(pos, regexp); matches != nil {
				if groupIndex >= len(matches) {
					panic(fmt.Sprintf("Capturing group %d is invalid for %s", groupIndex, regexp))
				}
				return ast.NewTerminalNode(schema, token, string(matches[groupIndex]), pos, readerPos), data.EmptyIntSet, nil
			}
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
