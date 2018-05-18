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

// Substring matches the given string
func Substring(token string, str string, value interface{}) *parser.NamedFunc {
	if str == "" {
		panic("Substring() should not be called with empty string")
	}

	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (parsley.Node, parsley.Error, data.IntSet) {
		tr := r.(*text.Reader)
		if readerPos, found := tr.MatchString(pos, str); found {
			return ast.NewTerminalNode(token, value, pos, readerPos), nil, data.EmptyIntSet
		}
		return nil, nil, data.EmptyIntSet
	}).WithName(fmt.Sprintf("%q", str))
}
