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

// Integer matches all integer numbers and zero with an optional -/+ sign
func Integer() parsley.ParserFunc {
	return parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
		tr := r.(*text.Reader)
		if readerPos, result := tr.ReadRegexp(pos, "[-+]?(?:[1-9][0-9]*|0[xX][0-9a-fA-F]+|0[0-7]*)"); result != nil {
			intValue, err := strconv.ParseInt(string(result), 0, 0)
			if err != nil {
				// This should never happen
				panic(fmt.Sprintf("Could not convert %s to integer", string(result)))
			}
			return data.EmptyIntSet, []parsley.Node{ast.NewTerminalNode("INT", int(intValue), r.Pos(pos), readerPos)}, nil
		}
		return data.EmptyIntSet, nil, parsley.NewError(r.Pos(pos), "was expecting integer value")
	})
}
