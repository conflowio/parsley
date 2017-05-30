package ast_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/stretchr/testify/assert"
)

func getInterpreterFunc(val interface{}, err error) ast.InterpreterFunc {
	return ast.InterpreterFunc(func(values []interface{}) (interface{}, error) {
		return val, err
	})
}

func assertNodesEqual(t *testing.T, expected ast.Node, actual ast.Node) {
	assert.Equal(t, expected.Token(), actual.Token())
	assert.Equal(t, expected.Pos(), actual.Pos())
	actualVal, actualErr := actual.Value()
	expectedVal, expectedErr := expected.Value()
	assert.Equal(t, expectedVal, actualVal)
	assert.Equal(t, expectedErr, actualErr)
}
