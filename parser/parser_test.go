package parser_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/parser"
	"github.com/opsidian/parsec/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectLeftRecursion(t *testing.T) {
	input := "abbbbbbbbbbbbbbbbbbbb"
	r := reader.New([]byte(input), true)
	var a parser.Func

	a = parser.Or("A",
		parser.And("AB", stringBuilder(),
			&a,
			parser.Rune('b', "CHAR"),
		),
		parser.Rune('a', "CHAR"),
	)
	all := parser.And("ALL", ast.SingleNodeBuilder(0), &a, parser.End())
	c := parser.NewContext()
	results := all.Parse(c, r)
	assert.Equal(t, 1, len(results.Items()))
	result, err := results.Items()[0].Node().Value()
	require.Nil(t, err)
	assert.Equal(t, input, result)
	assert.Equal(t, 300, c.GetSumCallCount())
}

func TestIndirectLeftRecursion(t *testing.T) {
	input := []byte("1 + 2 + 3")
	r := reader.New(input, true)
	var add parser.Func

	value := parser.Or(
		"VALUE",
		intLiteral(),
		&add,
	)

	add = parser.And(
		"ADD",
		ast.BinaryOperatorBuilder("ADD", func(children []interface{}) (interface{}, error) {
			return children[0].(int) + children[1].(int), nil
		}),
		value,
		parser.Rune('+', "ADD"),
		value,
	)
	p := parser.And("ALL", ast.SingleNodeBuilder(0), value, parser.End())
	c := parser.NewContext()
	results := p.Parse(c, r)
	result, err := results.Items()[0].Node().Value()
	require.Nil(t, err)
	assert.Equal(t, 6, result)
	assert.Equal(t, 1297, c.GetSumCallCount())
}

func stringBuilder() ast.NodeBuilder {
	return func(nodes []ast.Node) ast.Node {
		s := ""
		for _, node := range nodes {
			val, _ := node.Value()
			if runeVal, ok := val.(rune); ok {
				s = s + string(runeVal)
			} else {
				s = s + val.(string)
			}
		}
		return ast.NewTerminalNode("STRING", nodes[0].Pos(), s)
	}
}

func intLiteral() parser.Func {
	return parser.Func(func(c *parser.Context, r *reader.Reader) *parser.Results {
		if matches, pos := r.ReadMatch("^[\\-+]?[1-9][0-9]*"); matches != nil {
			intValue, err := strconv.Atoi(matches[0])
			if err != nil {
				panic(fmt.Sprintf("Could not convert %s to integer", matches[0]))
			}
			return parser.NewResult(ast.NewTerminalNode("INT", pos, intValue), r).AsList()
		}
		return nil
	})
}
