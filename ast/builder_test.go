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
	assertNodesEqual(t, expectedNode, actualNode)
}

func TestSingleNodeBuilderShouldReturnWithSelectednode(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewTerminalNode("2", test.NewPosition(1), "2"),
	}
	assert.Equal(t, nodes[1], ast.SingleNodeBuilder(1).BuildNode(nodes))
}

func TestSingleNodeBuilderShouldPanicIfIndexIsOutOfBounds(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewTerminalNode("2", test.NewPosition(1), "2"),
	}
	assert.Panics(t, func() { ast.SingleNodeBuilder(2).BuildNode(nodes) })
}

func TestAllNodesBuilderShouldIncludeAllNodes(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewTerminalNode("2", test.NewPosition(1), "2"),
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
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
		ast.NewTerminalNode("+", test.NewPosition(1), "+"),
		ast.NewTerminalNode("2", test.NewPosition(2), 2),
	}
	interpreter := getInterpreterFunc(3, nil)
	expected := ast.NewNonTerminalNode("+", []ast.Node{nodes[0], nodes[2]}, interpreter)
	actual := ast.BinaryOperatorBuilder(interpreter).BuildNode(nodes)
	assertNodesEqual(t, expected, actual)
}

func TestBinaryOperatorBuilderShouldPanicIfNotThreeNodes(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewTerminalNode("2", test.NewPosition(1), "2"),
	}
	assert.Panics(t, func() { ast.BinaryOperatorBuilder(nil).BuildNode(nodes) })
}

func TestNilBuilderReturnsWithNil(t *testing.T) {
	assert.Nil(t, ast.NilBuilder().BuildNode([]ast.Node{}))
}

func TestFlattenBuilderShouldReturnWithFlattenNodes(t *testing.T) {
	interpreter := getInterpreterFunc("X", nil)
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewNonTerminalNode(
			"3",
			[]ast.Node{
				ast.NewTerminalNode("2", test.NewPosition(1), "2"),
			},
			interpreter,
		),
	}
	expectedNodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewTerminalNode("2", test.NewPosition(1), "2"),
	}
	expected := ast.NewNonTerminalNode("TEST", expectedNodes, interpreter)
	actual := ast.FlattenBuilder("TEST", interpreter).BuildNode(nodes)
	assert.Equal(t, expected.Token(), actual.Token())
	assert.Equal(t, expected.Pos(), actual.Pos())
	actualVal, actualErr := actual.Value()
	expectedVal, expectedErr := expected.Value()
	assert.Equal(t, expectedVal, actualVal)
	assert.Equal(t, expectedErr, actualErr)
}
