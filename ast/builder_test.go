// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestNodeBuilderFuncShouldCallFunction(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewTerminalNode("2", test.NewPosition(1), "2"),
	}
	var actualNodes []ast.Node
	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
		return nil, nil
	})
	expectedNode := ast.NewNonTerminalNode("TEST", nodes, interpreter)
	nodeBuilderFunc := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		actualNodes = nodes
		return expectedNode
	})

	actualNode := nodeBuilderFunc.BuildNode(nodes)
	assert.Equal(t, nodes, actualNodes)
	test.AssertNodesEqual(t, nil, expectedNode, actualNode)
}
