// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"fmt"
	"strconv"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// Char matches a character literal enclosed in single quotes
func Char() parsley.ParserFunc {
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		tr := r.(*text.Reader)
		readerPos, res := tr.ReadRegexp(
			pos, "'(.|\\\\[abfnrtv]|\\\\x[0-9a-fA-F]{2,2}|\\\\u[0-9a-fA-F]{4,4}|\\\\U[0-9a-fA-F]{8,8})'",
		)
		if res != nil {
			value, _, tail, err := strconv.UnquoteChar(string(res[1:len(res)-1]), '"')
			if tail != "" {
				// This should never happen
				panic(fmt.Sprintf("Unprocessed string segment: %s", tail))
			}
			if err == nil {
				return data.EmptyIntSet(), []parsley.Node{ast.NewTerminalNode("CHAR", value, r.Pos(pos), readerPos)}, nil
			}
		}
		return data.EmptyIntSet(), nil, parsley.NewError(r.Pos(pos), "was expecting char literal")
	})
}
