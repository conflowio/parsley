package ast_test

import (
	"errors"
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestTerminalNode(t *testing.T) {
	pos := test.NewPosition(1)
	node := ast.NewTerminalNode("TOKEN", pos, "VALUE")
	assert.Equal(t, "TOKEN", node.Token())
	assert.Equal(t, pos, node.Pos())
	actualVal, actualErr := node.Value()
	assert.Equal(t, "VALUE", actualVal)
	assert.Nil(t, actualErr)

	assert.Equal(t, "T{VALUE, Pos{1}}", node.String())

	node = ast.NewTerminalNode("TOKEN", pos, nil)
	assert.Equal(t, "T{TOKEN, Pos{1}}", node.String())
}

func TestNonTerminalNode(t *testing.T) {
	expectedValue := 3
	expectedErr := errors.New("E")
	var actualNodes []ast.Node
	interpreterFunc := ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
		actualNodes = nodes
		return expectedValue, expectedErr
	})

	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
		ast.NewTerminalNode("2", test.NewPosition(2), 2),
	}

	node := ast.NewNonTerminalNode("+", nodes, interpreterFunc)
	assert.Equal(t, "+", node.Token())
	assert.Equal(t, nodes, node.Children())
	actualVal, actualErr := node.Value()
	assert.Equal(t, nodes, actualNodes)
	assert.Equal(t, expectedValue, actualVal)
	assert.Equal(t, expectedErr, actualErr)

	assert.Equal(t, "NT{+, [T{1, Pos{0}} T{2, Pos{2}}]}", node.String())
}

func TestNonTerminalNodeValueShouldIncludeNilNodes(t *testing.T) {
	var actualNodes []ast.Node
	interpreterFunc := ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
		actualNodes = nodes
		return nil, nil
	})

	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
		nil,
		ast.NewTerminalNode("2", test.NewPosition(2), 2),
	}

	node := ast.NewNonTerminalNode("+", nodes, interpreterFunc)
	node.Value()
	assert.Equal(t, nodes, actualNodes)
}

func TestNonTerminalNodeValueShouldReturnNilIfNoInterpreter(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("1", test.NewPosition(0), 1),
	}

	node := ast.NewNonTerminalNode("+", nodes, nil)
	val, err := node.Value()
	assert.Nil(t, val)
	assert.Nil(t, err)
}

func TestNonTerminalNodeValueShouldCallInterpreterEvenIfNoChildren(t *testing.T) {
	interpreterFunc := ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
		return 1, nil
	})
	node := ast.NewNonTerminalNode("+", []ast.Node{}, interpreterFunc)
	val, err := node.Value()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)
}

func TestNonTerminalNodeValueShouldReturnNilIfAllChildrenAreNil(t *testing.T) {
	interpreter := ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
		return nil, nil
	})
	nodes := []ast.Node{
		nil,
		nil,
	}
	node := ast.NewNonTerminalNode("+", nodes, interpreter)
	val, err := node.Value()
	assert.Nil(t, val)
	assert.Nil(t, err)
}
