package parsley_test

import (
	"errors"
	"testing"

	"github.com/opsidian/parsley"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseShouldRunParserAndReturnNode(t *testing.T) {
	expectedNode := ast.NewTerminalNode("STRING", text.NewPosition(1, 2, 3), "RES")
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), parser.NewResult(expectedNode, r).AsSet()
	})
	node, err := parsley.ParseText([]byte("input"), true, s)
	assert.Equal(t, expectedNode, node)
	assert.Nil(t, err)
}

func TestParseShouldHandleEmptyResult(t *testing.T) {
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})
	node, err := parsley.ParseText([]byte("input"), true, s)
	assert.Error(t, err)
	assert.Nil(t, node)
}

func TestEvaluateShouldRunParserAndReturnValue(t *testing.T) {
	expectedValue := "RES"
	node := ast.NewTerminalNode("STRING", text.NewPosition(1, 2, 3), expectedValue)
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), parser.NewResult(node, r).AsSet()
	})
	value, err := parsley.EvaluateText([]byte("input"), true, s)
	assert.Equal(t, expectedValue, value)
	assert.Nil(t, err)
}

func TestEvaluateShouldHandleEmptyResult(t *testing.T) {
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})
	value, err := parsley.EvaluateText([]byte("input"), true, s)
	assert.Error(t, err)
	assert.Nil(t, value)
}

func TestEvaluateShouldHandleInterpreterError(t *testing.T) {
	randomChild := ast.NewTerminalNode("X", text.NewPosition(1, 2, 3), "X")
	expectedErr := errors.New("ERR")
	node := ast.NewNonTerminalNode("ERR", []ast.Node{randomChild}, ast.InterpreterFunc(func([]ast.Node) (interface{}, error) {
		return nil, expectedErr
	}))
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), parser.NewResult(node, r).AsSet()
	})
	value, err := parsley.EvaluateText([]byte("input"), true, s)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, value)
}

func TestDirectLeftRecursion(t *testing.T) {
	input := "abbbbbbbbbbbbbbbbbbb"
	h := parser.NewHistory()

	var a parser.Func
	a = combinator.Memoize("A", h, combinator.Or(
		combinator.And(stringBuilder(),
			&a,
			terminal.Rune('b', "CHAR"),
		),
		terminal.Rune('a', "CHAR"),
	))
	s := combinator.And(builder.Select(0), &a, parser.End())

	result, err := parsley.EvaluateText([]byte(input), true, s)
	require.Nil(t, err)
	assert.Equal(t, input, result)
	assert.Equal(t, 318, parser.Stat.GetSumCallCount())
}

func TestIndirectLeftRecursion(t *testing.T) {
	input := []byte("1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10")
	h := parser.NewHistory()

	var add parser.Func
	value := combinator.Memoize("VALUE", h, combinator.Or(
		terminal.Integer(),
		&add,
	))

	add = combinator.Memoize("ADD", h, combinator.And(
		builder.BinaryOperation(
			ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
				value0, _ := nodes[0].Value()
				value1, _ := nodes[1].Value()
				return value0.(int) + value1.(int), nil
			}),
		),
		value,
		terminal.Rune('+', "ADD"),
		value,
	))
	s := combinator.And(builder.Select(0), value, parser.End())

	result, err := parsley.EvaluateText([]byte(input), true, s)
	require.Nil(t, err)
	assert.Equal(t, 55, result)
	assert.Equal(t, 3459, parser.Stat.GetSumCallCount())
}

func TestManySepBy(t *testing.T) {
	input := []byte("1 - 2 + 3 - 4 + 5 - 6 + 7 - 8 + 9 - 10")
	h := parser.NewHistory()

	var add parser.Func
	value := combinator.Memoize("VALUE", h, combinator.Or(
		terminal.Integer(),
		&add,
	))

	add = combinator.Memoize("SUM", h, combinator.ManySepBy(
		"SUM", "+", h, value, combinator.Or(terminal.Rune('+', "+"), terminal.Rune('-', "-")), 1,
		ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
			sum := 0
			modifier := 1
			for _, node := range nodes {
				v, err := node.Value()
				if err != nil {
					return nil, err
				}
				switch vt := v.(type) {
				case int:
					sum += modifier * vt
				case rune:
					if vt == '+' {
						modifier = 1
					} else {
						modifier = -1
					}
				}
			}
			return sum, nil
		}),
	))

	s := combinator.And(builder.Select(0), value, parser.End())
	result, err := parsley.EvaluateText([]byte(input), true, s)
	require.Nil(t, err)
	assert.Equal(t, -5, result)
	assert.Equal(t, 1246, parser.Stat.GetSumCallCount())
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
		first := nodes[0].(ast.TerminalNode)
		return ast.NewTerminalNode("STRING", first.Pos(), s)
	})
}
