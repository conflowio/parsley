package builder_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestSingleNodeBuilderShouldReturnWithSelectednode(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewTerminalNode("2", test.NewPosition(1), "2"),
	}
	assert.Equal(t, nodes[1], builder.Select(1).BuildNode(nodes))
}

func TestSingleNodeBuilderShouldPanicIfIndexIsOutOfBounds(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewTerminalNode("2", test.NewPosition(1), "2"),
	}
	assert.Panics(t, func() { builder.Select(2).BuildNode(nodes) })
}

func TestAllNodesBuilderShouldIncludeAllNodes(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewTerminalNode("2", test.NewPosition(1), "2"),
	}
	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
		return nil, nil
	})
	expected := ast.NewNonTerminalNode("TEST", nodes, interpreter)
	actual := builder.All("TEST", interpreter).BuildNode(nodes)
	test.AssertNodesEqual(t, nil, expected, actual)
}

func TestBinaryOperatorBuilderShouldBuildNode(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
		ast.NewTerminalNode("+", test.NewPosition(1), "+"),
		ast.NewTerminalNode("2", test.NewPosition(2), 2),
	}
	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
		return nil, nil
	})
	expected := ast.NewNonTerminalNode("+", []ast.Node{nodes[0], nodes[2]}, interpreter)
	actual := builder.BinaryOperation(interpreter).BuildNode(nodes)
	test.AssertNodesEqual(t, nil, expected, actual)
}

func TestBinaryOperatorBuilderShouldPanicIfNotThreeNodes(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), "1"),
		ast.NewTerminalNode("2", test.NewPosition(1), "2"),
	}
	assert.Panics(t, func() { builder.BinaryOperation(nil).BuildNode(nodes) })
}

func TestNilBuilderReturnsWithNil(t *testing.T) {
	assert.Nil(t, builder.Nil().BuildNode([]ast.Node{}))
}

func TestFlattenBuilderShouldReturnWithFlattenNodes(t *testing.T) {
	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
		return nil, nil
	})
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
	actual := builder.Flatten("TEST", interpreter).BuildNode(nodes)
	test.AssertNodesEqual(t, nil, expected, actual)
}
