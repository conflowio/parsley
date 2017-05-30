package ast_test

import (
	"errors"
	"testing"

	"github.com/opsidian/parsley/ast"
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
