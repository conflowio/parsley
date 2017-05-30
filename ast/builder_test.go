package ast_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
)

func TestNodeBuilderFuncShouldCallFunction(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", reader.NewPosition(0, 1, 1), "1"),
		ast.NewTerminalNode("2", reader.NewPosition(1, 1, 2), "2"),
	}
	var actualNodes []ast.Node
	expectedNode := ast.NewNonTerminalNode("TEST", nodes, getInterpreterFunc("X", nil))
	nodeBuilderFunc := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		actualNodes = nodes
		return expectedNode
	})

	actualNode := nodeBuilderFunc.BuildNode(nodes)
	assert.Equal(t, nodes, actualNodes)
	assertNodesEqual(t, expectedNode, actualNode)
}

func TestSingleNodeBuilderShouldReturnWithSelectednode(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", reader.NewPosition(0, 1, 1), "1"),
		ast.NewTerminalNode("2", reader.NewPosition(1, 1, 2), "2"),
	}
	assert.Equal(t, nodes[1], ast.SingleNodeBuilder(1).BuildNode(nodes))
}

func TestSingleNodeBuilderShouldPanicIfIndexIsOutOfBounds(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", reader.NewPosition(0, 1, 1), "1"),
		ast.NewTerminalNode("2", reader.NewPosition(1, 1, 2), "2"),
	}
	assert.Panics(t, func() { ast.SingleNodeBuilder(2).BuildNode(nodes) })
}

func TestAllNodesBuilderShouldIncludeAllNodes(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", reader.NewPosition(0, 1, 1), "1"),
		ast.NewTerminalNode("2", reader.NewPosition(1, 1, 2), "2"),
	}
	interpreter := getInterpreterFunc("X", nil)
	expected := ast.NewNonTerminalNode("TEST", nodes, interpreter)
	actual := ast.AllNodesBuilder("TEST", interpreter).BuildNode(nodes)
	assert.Equal(t, expected.Token(), actual.Token())
	assert.Equal(t, expected.Pos(), actual.Pos())
	actualVal, actualErr := actual.Value()
	expectedVal, expectedErr := expected.Value()
	assert.Equal(t, expectedVal, actualVal)
	assert.Equal(t, expectedErr, actualErr)
}

func TestBinaryOperatorBuilderShouldBuildNode(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", reader.NewPosition(0, 1, 1), 1),
		ast.NewTerminalNode("+", reader.NewPosition(1, 1, 2), "+"),
		ast.NewTerminalNode("2", reader.NewPosition(2, 1, 3), 2),
	}
	interpreter := getInterpreterFunc(3, nil)
	expected := ast.NewNonTerminalNode("+", []ast.Node{nodes[0], nodes[2]}, interpreter)
	actual := ast.BinaryOperatorBuilder(interpreter).BuildNode(nodes)
	assertNodesEqual(t, expected, actual)
}

func TestBinaryOperatorBuilderShouldPanicIfNotThreeNodes(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", reader.NewPosition(0, 1, 1), "1"),
		ast.NewTerminalNode("2", reader.NewPosition(1, 1, 2), "2"),
	}
	assert.Panics(t, func() { ast.BinaryOperatorBuilder(nil).BuildNode(nodes) })
}
