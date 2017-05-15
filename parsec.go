package main

import (
	"fmt"
	"go/token"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/parser"
	"github.com/opsidian/parsec/reader"
)

func main() {
	input := []byte("1 + 2 - 3 + 4 - 5 + 6 - 7 + 8 - 9 + 10 - 11 + 12 - 13")
	r := reader.New(input)
	var expr parser.Func

	value := parser.Or(
		"VALUE",
		parser.IntLiteral(),
		&expr,
	)

	add := parser.And(
		"ADD",
		ast.BinaryOperatorBuilder(token.ADD, func(children []interface{}) (interface{}, error) {
			return children[0].(int) + children[1].(int), nil
		}),
		value,
		parser.Rune('+', token.ADD),
		value,
	)
	subtract := parser.And(
		"SUB",
		ast.BinaryOperatorBuilder(token.SUB, func(children []interface{}) (interface{}, error) {
			return children[0].(int) - children[1].(int), nil
		}),
		value,
		parser.Rune('-', token.SUB),
		value,
	)
	expr = parser.Or("EXPR", add, subtract)
	all := parser.And("ALL", ast.SingleNodeBuilder(0), value, parser.End())
	h := parser.NewHistory()
	results := all.Parse(h, r)
	if results == nil {
		panic(fmt.Sprintf("Couldn't parse the expression: %s", input))
	}
	result, err := results[0].Node().Value()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result was: %v, calls: %d, number of AST trees: %d\n", result, h.GetCallCount(), len(results))
}
