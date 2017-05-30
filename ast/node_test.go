package ast_test

import (
	"errors"
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
)

func TesInterpreterFuncShouldCallFunction(t *testing.T) {
	values := []interface{}{1, 2}
	expectedValue := 1
	expectedErr := errors.New("e")
	var actualValues []interface{}
	interpreterFunc := ast.InterpreterFunc(func(values []interface{}) (interface{}, error) {
		actualValues = values
		return expectedValue, expectedErr
	})

	actualValue, actualErr := interpreterFunc.Eval(values)

	assert.Equal(t, values, actualValues)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, expectedErr, actualErr)
}

func TestTerminalNode(t *testing.T) {
	pos := reader.NewPosition(1, 2, 3)
	node := ast.NewTerminalNode("TOKEN", pos, "VALUE")
	assert.Equal(t, "TOKEN", node.Token())
	assert.Equal(t, pos, node.Pos())
	actualVal, actualErr := node.Value()
	assert.Equal(t, "VALUE", actualVal)
	assert.Nil(t, actualErr)

	assert.Equal(t, "T{VALUE, 2:3}", node.String())

	node = ast.NewTerminalNode("TOKEN", pos, nil)
	assert.Equal(t, "T{TOKEN, 2:3}", node.String())
}

func TestNonTerminalNode(t *testing.T) {
	expectedValue := 3
	expectedErr := errors.New("E")
	var actualValues []interface{}
	interpreterFunc := ast.InterpreterFunc(func(values []interface{}) (interface{}, error) {
		actualValues = values
		return expectedValue, expectedErr
	})

	nodes := []ast.Node{
		ast.NewTerminalNode("1", reader.NewPosition(0, 1, 2), 1),
		ast.NewTerminalNode("2", reader.NewPosition(2, 1, 4), 2),
	}

	node := ast.NewNonTerminalNode("+", nodes, interpreterFunc)
	assert.Equal(t, "+", node.Token())
	assert.Equal(t, nodes[0].Pos(), node.Pos())
	actualVal, actualErr := node.Value()
	assert.Equal(t, []interface{}{1, 2}, actualValues)
	assert.Equal(t, expectedValue, actualVal)
	assert.Equal(t, expectedErr, actualErr)

	assert.Equal(t, "NT{+, [T{1, 1:2} T{2, 1:4}]}", node.String())
}

func TestNonTerminalNodeShouldPanicWithoutChildren(t *testing.T) {
	assert.Panics(t, func() { ast.NewNonTerminalNode("X", nil, getInterpreterFunc(nil, nil)) })
}

func TestNonTerminalNodeShouldPanicWithoutInterpreter(t *testing.T) {
	randomNode := ast.NewTerminalNode("X", reader.NewPosition(1, 2, 3), getInterpreterFunc(nil, nil))
	assert.Panics(t, func() { ast.NewNonTerminalNode("X", []ast.Node{randomNode}, nil) })
}

func TestNonTerminalNodeValueShouldReturnErrorIfChildHasError(t *testing.T) {
	expectedErr := errors.New("E")
	randomNode := ast.NewTerminalNode("X", reader.NewPosition(1, 2, 3), getInterpreterFunc(nil, nil))
	badChild := ast.NewNonTerminalNode("BAD", []ast.Node{randomNode}, ast.InterpreterFunc(func(values []interface{}) (interface{}, error) {
		return nil, expectedErr
	}))
	node := ast.NewNonTerminalNode("X", []ast.Node{badChild}, getInterpreterFunc(nil, nil))
	actualVal, actualErr := node.Value()
	assert.Nil(t, actualVal)
	assert.Equal(t, expectedErr, actualErr)
}
