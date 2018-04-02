// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package interpreter

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"
)

// Select returns with an interpreter function which returns the value of the selected node
func Select(i int) ast.InterpreterFunc {
	return func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		if i >= len(nodes) {
			panic(fmt.Sprintf("node index is out of bounds: %d", i))
		}
		return nodes[i].Value(ctx)
	}
}

// Nil returns with an interpreter function which always returns with a nil result
func Nil() ast.InterpreterFunc {
	return func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		return nil, nil
	}
}
