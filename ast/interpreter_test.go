package ast_test

import (
	"errors"
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TesInterpreterFuncShouldCallFunction(t *testing.T) {
	nodes := []ast.Node{
		ast.NewTerminalNode("INT", test.NewPosition(0), 1),
	}
	expectedValue := 1
	expectedErr := errors.New("e")
	var actualNodes []ast.Node
	interpreterFunc := ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
		actualNodes = nodes
		return expectedValue, expectedErr
	})

	actualValue, actualErr := interpreterFunc.Eval(nodes)

	assert.Equal(t, nodes, actualNodes)
	assert.Equal(t, expectedValue, actualValue)
	assert.Equal(t, expectedErr, actualErr)
}
