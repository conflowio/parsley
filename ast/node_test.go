// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/reader/mocks"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestTerminalNode(t *testing.T) {
	pos := test.NewPosition(1)
	node := ast.NewTerminalNode("TOKEN", pos, "VALUE")
	assert.Equal(t, "TOKEN", node.Token())
	assert.Equal(t, pos, node.Pos())
	actualVal, actualErr := node.Value(nil)
	assert.Equal(t, "VALUE", actualVal)
	assert.Nil(t, actualErr)

	assert.Equal(t, "T{VALUE, Pos{1}}", node.String())

	node = ast.NewTerminalNode("TOKEN", pos, nil)
	assert.Equal(t, "T{TOKEN, Pos{1}}", node.String())
}

func TestNonTerminalNode(t *testing.T) {
	expectedValue := 3
	expectedErr := new(mocks.Error)
	var actualCtx interface{}
	var actualNodes []ast.Node
	interpreterFunc := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		actualCtx = ctx
		actualNodes = nodes
		return expectedValue, expectedErr
	})

	ctx := "textctx"
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
		ast.NewTerminalNode("2", test.NewPosition(2), 2),
	}

	node := ast.NewNonTerminalNode("+", nodes, interpreterFunc)
	assert.Equal(t, "+", node.Token())
	assert.Equal(t, nodes, node.Children())
	assert.Equal(t, test.NewPosition(0), node.Pos())
	actualVal, actualErr := node.Value(ctx)
	assert.Equal(t, ctx, actualCtx)
	assert.Equal(t, nodes, actualNodes)
	assert.Equal(t, expectedValue, actualVal)
	assert.Equal(t, expectedErr, actualErr)

	assert.Equal(t, "NT{+, [T{1, Pos{0}} T{2, Pos{2}}]}", node.String())
}

func TestNonTerminalNodeShouldGetPosFromFirstNonNilChild(t *testing.T) {
	nodes := []ast.Node{
		nil,
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
	}

	node := ast.NewNonTerminalNode("TEST", nodes, nil)
	assert.Equal(t, test.NewPosition(0), node.Pos())
}

func TestNonTerminalNodeValueShouldIncludeNilNodes(t *testing.T) {
	var actualNodes []ast.Node
	interpreterFunc := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		actualNodes = nodes
		return nil, nil
	})

	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
		nil,
		ast.NewTerminalNode("2", test.NewPosition(2), 2),
	}

	node := ast.NewNonTerminalNode("+", nodes, interpreterFunc)
	node.Value(nil)
	assert.Equal(t, nodes, actualNodes)
}

func TestNonTerminalNodeValueShouldReturnNilIfNoInterpreter(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
	}

	node := ast.NewNonTerminalNode("+", nodes, nil)
	val, err := node.Value(nil)
	assert.Nil(t, val)
	assert.Nil(t, err)
}

func TestNonTerminalNodeValueShouldCallInterpreterEvenIfNoChildren(t *testing.T) {
	interpreterFunc := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		return 1, nil
	})
	node := ast.NewNonTerminalNode("+", []ast.Node{}, interpreterFunc)
	val, err := node.Value(nil)
	assert.Equal(t, 1, val)
	assert.Nil(t, err)
}

func TestNonTerminalNodeValueShouldReturnNilIfAllChildrenAreNil(t *testing.T) {
	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		return nil, nil
	})
	nodes := []ast.Node{
		nil,
		nil,
	}
	node := ast.NewNonTerminalNode("+", nodes, interpreter)
	val, err := node.Value(nil)
	assert.Nil(t, val)
	assert.Nil(t, err)
}
