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
		if i < 0 || i >= len(nodes) {
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

// Array can be used to create an array from a list of nodes, where values and separators are following
// each-other
func Array() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		res := make([]interface{}, (len(nodes)+1)/2)
		for i := 0; i < len(nodes); i += 2 {
			value, err := nodes[i].Value(ctx)
			if err != nil {
				return nil, err
			}
			res[i/2] = value
		}
		return res, nil
	})
}

// Object can be used to create an object from a list of nodes, where key-value nodes and separater nodes
// follow each-other, and a key-value node consists of a key node, a separator node and a value node
func Object() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		res := make(map[string]interface{}, (len(nodes)+1)/2)
		for i := 0; i < len(nodes); i += 2 {
			keyValue := nodes[i].(*ast.NonTerminalNode)
			key, err := keyValue.Children()[0].Value(ctx)
			if err != nil {
				return nil, err
			}
			value, err := keyValue.Children()[2].Value(ctx)
			if err != nil {
				return nil, err
			}
			res[key.(string)] = value
		}
		return res, nil
	})
}
