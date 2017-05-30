package parser_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectLeftRecursion(t *testing.T) {
	input := "abbbbbbbbbbbbbbbbbbb"
	r := reader.New([]byte(input), true)
	h := parser.NewHistory()
	var a parser.Func

	a = parser.Or("A", h,
		parser.And("AB", h, stringBuilder(),
			&a,
			parser.Rune('b', "CHAR"),
		),
		parser.Rune('a', "CHAR"),
	)
	all := parser.And("ALL", h, ast.SingleNodeBuilder(0), &a, parser.End())
	results := all.Parse(data.NewIntMap(nil), r)
	assert.Equal(t, 1, len(results.Results))
	result, err := results.Results[0].Node().Value()
	require.Nil(t, err)
	assert.Equal(t, input, result)
	assert.Equal(t, 318, h.GetSumCallCount())
}

func TestIndirectLeftRecursion(t *testing.T) {
	input := []byte("1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10")
	r := reader.New(input, true)
	h := parser.NewHistory()
	var add parser.Func

	value := parser.Or("VALUE", h,
		intLiteral(),
		&add,
	)

	add = parser.And("ADD", h,
		ast.BinaryOperatorBuilder(
			ast.InterpreterFunc(func(children []interface{}) (interface{}, error) {
				return children[0].(int) + children[1].(int), nil
			}),
		),
		value,
		parser.Rune('+', "ADD"),
		value,
	)
	p := parser.And("ALL", h, ast.SingleNodeBuilder(0), value, parser.End())
	results := p.Parse(data.NewIntMap(nil), r)
	require.Equal(t, 1, len(results.Results), "Parser should be successful")
	result, err := results.Results[0].Node().Value()
	require.Nil(t, err)
	assert.Equal(t, 55, result)
	assert.Equal(t, 3459, h.GetSumCallCount())
}

func TestMany(t *testing.T) {
	input := []byte("1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10")
	r := reader.New(input, true)
	h := parser.NewHistory()
	var add parser.Func

	value := parser.Or("VALUE", h,
		intLiteral(),
		&add,
	)

	add = parser.ManySep(
		"ADD", "+", h, value, parser.Rune('+', "+"),
		ast.InterpreterFunc(func(children []interface{}) (interface{}, error) {
			sum := 0
			for _, v := range children {
				sum += v.(int)
			}
			return sum, nil
		}),
	)

	p := parser.And("ALL", h, ast.SingleNodeBuilder(0), value, parser.End())
	results := p.Parse(data.NewIntMap(nil), r)
	require.Equal(t, 1, len(results.Results), "Parser should be successful")
	result, err := results.Results[0].Node().Value()
	require.Nil(t, err)
	assert.Equal(t, 55, result)
	assert.Equal(t, 1226, h.GetSumCallCount())
}

func stringBuilder() ast.NodeBuilder {
	return ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
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
	})
}

func intLiteral() parser.Func {
	return parser.Func(func(ctx data.IntMap, r *reader.Reader) *parser.ParserResult {
		if matches, pos := r.ReadMatch("^0|[\\-+]?[1-9][0-9]*"); matches != nil {
			intValue, err := strconv.Atoi(matches[0])
			if err != nil {
				panic(fmt.Sprintf("Could not convert %s to integer", matches[0]))
			}
			return parser.NewParserResult(data.NewIntSet(), parser.NewResult(ast.NewTerminalNode("INT", pos, intValue), r))
		}
		return nil
	})
}
