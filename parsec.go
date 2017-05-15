package main

import (
	"fmt"
	"go/token"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/parser"
	"github.com/opsidian/parsec/reader"
)

func main() {
	input := []byte("(1 + (2 - 3))")
	r := reader.New(input)
	var expr parser.Func

	value := parser.Or(
		parser.IntLiteral(),
		parser.And(
			ast.SingleNodeBuilder(1),
			parser.Rune('(', token.LPAREN),
			&expr,
			parser.Rune(')', token.RPAREN),
		),
	)

	add := parser.And(
		ast.BinaryOperatorBuilder(token.ADD, func(children []interface{}) (interface{}, error) {
			return children[0].(int) + children[1].(int), nil
		}),
		value,
		parser.Rune('+', token.ADD),
		value,
	)
	subtract := parser.And(
		ast.BinaryOperatorBuilder(token.SUB, func(children []interface{}) (interface{}, error) {
			return children[0].(int) - children[1].(int), nil
		}),
		value,
		parser.Rune('-', token.SUB),
		value,
	)
	expr = parser.Or(add, subtract)
	all := parser.And(ast.SingleNodeBuilder(0), value, parser.End())
	results := all.Parse(r)
	if results == nil {
		panic(fmt.Sprintf("Couldn't parse the expression: %s", input))
	}
	result, err := results[0].Node().Value()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result was: %v\n", result)
}
