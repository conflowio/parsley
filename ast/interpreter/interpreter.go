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

// Select returns with an interpreter which returns the value of the selected node
func Select(i int) parsley.Interpreter {
	return selectInterpreter{i: i}
}

type selectInterpreter struct {
	i int
}

// StaticCheck runs the static checking on the indexed node
func (s selectInterpreter) StaticCheck(userCtx interface{}, node parsley.NonTerminalNode) (string, parsley.Error) {
	nodes := node.Children()
	if s.i < 0 || s.i >= len(nodes) {
		panic(fmt.Sprintf("node index is out of bounds: %d", s.i))
	}
	if staticCheckableNode, ok := nodes[s.i].(parsley.StaticCheckable); ok {
		if err := staticCheckableNode.StaticCheck(userCtx); err != nil {
			return "", err
		}
		return nodes[s.i].Type(), nil
	}
	return "", nil
}

func (s selectInterpreter) Eval(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
	nodes := node.Children()
	if s.i < 0 || s.i >= len(nodes) {
		panic(fmt.Sprintf("node index is out of bounds: %d", s.i))
	}
	return nodes[s.i].Value(userCtx)
}

// Nil returns with an interpreter function which always returns with a nil result
func Nil() ast.InterpreterFunc {
	return func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
		return nil, nil
	}
}

// Array can be used to create an array from a list of nodes, where values and separators are following
// each-other
func Array() ast.InterpreterFunc {
	return ast.InterpreterFunc(func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
		nodes := node.Children()
		res := make([]interface{}, (len(nodes)+1)/2)
		for i := 0; i < len(nodes); i += 2 {
			value, err := nodes[i].Value(userCtx)
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
	return ast.InterpreterFunc(func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
		nodes := node.Children()
		res := make(map[string]interface{}, (len(nodes)+1)/2)
		for i := 0; i < len(nodes); i += 2 {
			keyValue := nodes[i].(parsley.NonTerminalNode)
			key, err := keyValue.Children()[0].Value(userCtx)
			if err != nil {
				return nil, err
			}
			value, err := keyValue.Children()[2].Value(userCtx)
			if err != nil {
				return nil, err
			}
			res[key.(string)] = value
		}
		return res, nil
	})
}
