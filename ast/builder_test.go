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
	expectedNode := ast.NewNonTerminalNode("TEST", nodes, getInterpreterFunc("X", nil))
	nodeBuilderFunc := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		actualNodes = nodes
		return expectedNode
	})

	actualNode := nodeBuilderFunc.BuildNode(nodes)
	assert.Equal(t, nodes, actualNodes)
	test.AssertNodesEqual(t, expectedNode, actualNode)
}
