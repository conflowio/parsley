package parsley_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/opsidian/parsley"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseShouldRunParserAndReturnNode(t *testing.T) {
	expectedNode := ast.NewTerminalNode("STRING", reader.NewPosition(1, 2, 3), "RES")
	s := parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		return parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(expectedNode, r))
	})
	node, err := parsley.Parse([]byte("input"), true, s)
	assert.Equal(t, expectedNode, node)
	assert.Nil(t, err)
}

func TestParseShouldHandleEmptyResult(t *testing.T) {
	s := parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		return parser.NewParserResult(parser.NoCurtailingParsers())
	})
	node, err := parsley.Parse([]byte("input"), true, s)
	assert.Error(t, err)
	assert.Nil(t, node)
}

func TestEvaluateShouldRunParserAndReturnValue(t *testing.T) {
	expectedValue := "RES"
	node := ast.NewTerminalNode("STRING", reader.NewPosition(1, 2, 3), expectedValue)
	s := parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		return parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(node, r))
	})
	value, err := parsley.Evaluate([]byte("input"), true, s)
	assert.Equal(t, expectedValue, value)
	assert.Nil(t, err)
}

func TestEvaluateShouldHandleEmptyResult(t *testing.T) {
	s := parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		return parser.NewParserResult(parser.NoCurtailingParsers())
	})
	value, err := parsley.Evaluate([]byte("input"), true, s)
	assert.Error(t, err)
	assert.Nil(t, value)
}

func TestEvaluateShouldHandleInterpreterError(t *testing.T) {
	randomChild := ast.NewTerminalNode("X", reader.NewPosition(1, 2, 3), "X")
	expectedErr := errors.New("ERR")
	node := ast.NewNonTerminalNode("ERR", []ast.Node{randomChild}, ast.InterpreterFunc(func([]interface{}) (interface{}, error) {
		return nil, expectedErr
	}))
	s := parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		return parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(node, r))
	})
	value, err := parsley.Evaluate([]byte("input"), true, s)
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
	s := combinator.And(ast.SingleNodeBuilder(0), &a, terminal.End())

	result, err := parsley.Evaluate([]byte(input), true, s)
	require.Nil(t, err)
	assert.Equal(t, input, result)
	assert.Equal(t, 318, parser.Stat.GetSumCallCount())
}

func TestIndirectLeftRecursion(t *testing.T) {
	input := []byte("1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10")
	h := parser.NewHistory()

	var add parser.Func
	value := combinator.Memoize("VALUE", h, combinator.Or(
		intLiteral(),
		&add,
	))

	add = combinator.Memoize("ADD", h, combinator.And(
		ast.BinaryOperatorBuilder(
			ast.InterpreterFunc(func(children []interface{}) (interface{}, error) {
				return children[0].(int) + children[1].(int), nil
			}),
		),
		value,
		terminal.Rune('+', "ADD"),
		value,
	))
	s := combinator.And(ast.SingleNodeBuilder(0), value, terminal.End())

	result, err := parsley.Evaluate([]byte(input), true, s)
	require.Nil(t, err)
	assert.Equal(t, 55, result)
	assert.Equal(t, 3459, parser.Stat.GetSumCallCount())
}

func TestMany(t *testing.T) {
	input := []byte("1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10")
	h := parser.NewHistory()

	var add parser.Func
	value := combinator.Memoize("VALUE", h, combinator.Or(
		intLiteral(),
		&add,
	))

	add = combinator.Memoize("ADD", h, combinator.SepBy(
		"ADD", "+", h, value, terminal.Rune('+', "+"),
		ast.InterpreterFunc(func(children []interface{}) (interface{}, error) {
			sum := 0
			for _, v := range children {
				sum += v.(int)
			}
			return sum, nil
		}),
	))

	s := combinator.And(ast.SingleNodeBuilder(0), value, terminal.End())
	result, err := parsley.Evaluate([]byte(input), true, s)
	require.Nil(t, err)
	assert.Equal(t, 55, result)
	assert.Equal(t, 1226, parser.Stat.GetSumCallCount())
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
			return parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("INT", pos, intValue), r))
		}
		return nil
	})
}
