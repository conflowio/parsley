package main

import (
	"fmt"
	"go/token"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/parser"
	"github.com/opsidian/parsec/reader"
)

func main() {
	directLeftRec()
}

func indirectLeftRec() {
	input := []byte("((1 - 2) + (3 - 4)) + ((5 - 6) + (7 - 8) + (9 - 10))")
	r := reader.New(input, true)
	var expr, valueParen parser.Func

	value := parser.Or(
		"VALUE",
		parser.IntLiteral(),
		&expr,
		&valueParen,
	)

	valueParen = parser.And(
		"VALUE_PAREN",
		ast.SingleNodeBuilder(1),
		parser.Rune('(', token.LPAREN),
		&value,
		parser.Rune(')', token.RPAREN),
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
	expr = parser.Or("EXPR", subtract, add)
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
	fmt.Printf("AST: %s\n", results[0].Node())
}

func directLeftRec() {
	input := []byte("abbbbbbbbb")
	fmt.Printf("%s\n\n", string(input))
	r := reader.New(input, true)
	var a parser.Func

	a = parser.Or("A",
		parser.And("AB", allNodesBuilder(),
			&a,
			parser.Rune('b', token.CHAR),
		),
		parser.Rune('a', token.CHAR),
	)
	all := parser.And("ALL", ast.SingleNodeBuilder(0), &a, parser.End())
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
	fmt.Printf("AST: %s\n", results[0].Node())
}

func allNodesBuilder() ast.NodeBuilder {
	return func(nodes []ast.Node) ast.Node {
		return ast.NewNonTerminalNode(token.ADD, nodes, func(children []interface{}) (interface{}, error) {
			s := ""
			for _, child := range children {
				s = s + string(child.(string))
			}
			return s, nil
		})
	}
}
