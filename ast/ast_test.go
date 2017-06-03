package ast_test

import (
	"github.com/opsidian/parsley/ast"
)

func getInterpreterFunc(val interface{}, err error) ast.InterpreterFunc {
	return ast.InterpreterFunc(func(values []interface{}) (interface{}, error) {
		return val, err
	})
}
