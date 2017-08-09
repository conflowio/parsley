// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley_test

import (
	"errors"
	"fmt"
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

// Let's define a simple parser which is able to parse adding two integers.
func ExampleParseText() {
	add := combinator.Seq(
		builder.BinaryOperation(
			ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
				value0, _ := nodes[0].Value(ctx)
				value1, _ := nodes[1].Value(ctx)
				return value0.(int) + value1.(int), nil
			}),
		),
		terminal.Integer(),
		terminal.Rune('+', "ADD"),
		terminal.Integer(),
	)
	s := combinator.Seq(builder.Select(0), add, parser.End())

	node, err := parsley.ParseText([]byte("1 + 2"), true, s)
	if err != nil {
		panic(err)
	}

	value, err2 := node.Value(nil)
	if err2 != nil {
		panic(err2)
	}
	fmt.Printf("Result: %d\n", value.(int))
	// Output: Result: 3
}

// Let's define a simple parser which is able to parse adding two integers.
func ExampleEvaluateText() {
	add := combinator.Seq(
		builder.BinaryOperation(
			ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
				value0, _ := nodes[0].Value(ctx)
				value1, _ := nodes[1].Value(ctx)
				return value0.(int) + value1.(int), nil
			}),
		),
		terminal.Integer(),
		terminal.Rune('+', "ADD"),
		terminal.Integer(),
	)
	s := combinator.Seq(builder.Select(0), add, parser.End())

	value, err := parsley.EvaluateText([]byte("1 + 2"), true, s, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %d\n", value.(int))
	// Output: Result: 3
}

func TestParseShouldRunParserAndReturnNode(t *testing.T) {
	expectedNode := ast.NewTerminalNode("STRING", text.NewPosition(1, 2, 3), "RES")
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), parser.NewResult(expectedNode, r).AsSet(), nil
	})
	node, err := parsley.ParseText([]byte("input"), true, s)
	assert.Equal(t, expectedNode, node)
	assert.Nil(t, err)
}

func TestParseShouldHandleEmptyResult(t *testing.T) {
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(text.NewPosition(2, 1, 3), "encountered a test error")
	})
	node, err := parsley.ParseText([]byte("input"), true, s)
	assert.Error(t, err)
	assert.Equal(t, "Failed to parse the input, encountered a test error at 1:3", err.Error())
	assert.Nil(t, node)
}

func TestParseShouldHandleNilNode(t *testing.T) {
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), parser.NewResult(nil, r).AsSet(), nil
	})
	node, err := parsley.ParseText([]byte(""), true, s)
	assert.Nil(t, err)
	assert.Nil(t, node)
}

func TestEvaluateShouldRunParserAndReturnValue(t *testing.T) {
	expectedValue := "RES"
	node := ast.NewTerminalNode("STRING", text.NewPosition(1, 2, 3), expectedValue)
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), parser.NewResult(node, r).AsSet(), nil
	})
	value, err := parsley.EvaluateText([]byte("input"), true, s, nil)
	assert.Equal(t, expectedValue, value)
	assert.Nil(t, err)
}

func TestEvaluateShouldPassContext(t *testing.T) {
	ctx := "testctx"
	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
		return ctx, nil
	})
	node := ast.NewNonTerminalNode("STRING", nil, interpreter)
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), parser.NewResult(node, r).AsSet(), nil
	})
	value, err := parsley.EvaluateText([]byte("input"), true, s, ctx)
	assert.Equal(t, ctx, value)
	assert.Nil(t, err)
}

func TestEvaluateShouldHandleEmptyResult(t *testing.T) {
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(text.NewPosition(2, 1, 3), "encountered a test error")
	})
	value, err := parsley.EvaluateText([]byte("input"), true, s, nil)
	assert.Error(t, err)
	assert.Equal(t, "Failed to parse the input, encountered a test error at 1:3", err.Error())
	assert.Nil(t, value)
}

func TestEvaluateShouldHandleNilNode(t *testing.T) {
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), parser.NewResult(nil, r).AsSet(), nil
	})
	value, err := parsley.EvaluateText([]byte(""), true, s, nil)
	assert.Nil(t, err)
	assert.Nil(t, value)
}

func TestEvaluateShouldHandleInterpreterError(t *testing.T) {
	randomChild := ast.NewTerminalNode("X", text.NewPosition(1, 2, 3), "X")
	expectedErr := errors.New("ERR")
	node := ast.NewNonTerminalNode("ERR", []ast.Node{randomChild}, ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
		return nil, expectedErr
	}))
	s := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), parser.NewResult(node, r).AsSet(), nil
	})
	value, err := parsley.EvaluateText([]byte("input"), true, s, nil)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, value)
}

func TestDirectLeftRecursion(t *testing.T) {
	input := "abbbbbbbbbbbbbbbbbbb"
	h := parser.NewHistory()

	concatNodes := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
		s := ""
		for _, node := range nodes {
			val, err := node.Value(ctx)
			if err != nil {
				return nil, err
			}
			if runeVal, ok := val.(rune); ok {
				s = s + string(runeVal)
			} else {
				s = s + val.(string)
			}
		}
		return s, nil
	})

	var a parser.Func
	a = combinator.Memoize("A", h, combinator.Any("a or ab",
		combinator.Seq(builder.All("AB", concatNodes),
			&a,
			terminal.Rune('b', "CHAR"),
		),
		terminal.Rune('a', "CHAR"),
	))
	s := combinator.Seq(builder.Select(0), &a, parser.End())

	result, err := parsley.EvaluateText([]byte(input), true, s, nil)
	require.Nil(t, err)
	assert.Equal(t, input, result)
	assert.Equal(t, 318, parser.Stat.GetSumCallCount())
}

func TestIndirectLeftRecursion(t *testing.T) {
	input := []byte("1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10")
	h := parser.NewHistory()

	var add parser.Func
	value := combinator.Memoize("VALUE", h, combinator.Any("value",
		terminal.Integer(),
		&add,
	))

	add = combinator.Memoize("ADD", h, combinator.Seq(
		builder.BinaryOperation(
			ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
				value0, _ := nodes[0].Value(ctx)
				value1, _ := nodes[1].Value(ctx)
				return value0.(int) + value1.(int), nil
			}),
		),
		value,
		terminal.Rune('+', "ADD"),
		value,
	))
	s := combinator.Seq(builder.Select(0), value, parser.End())

	result, err := parsley.EvaluateText([]byte(input), true, s, nil)
	require.Nil(t, err)
	assert.Equal(t, 55, result)
	assert.Equal(t, 3459, parser.Stat.GetSumCallCount())
}

func TestSepBy(t *testing.T) {
	input := []byte("1 - 2 + 3 - 4 + 5 - 6 + 7 - 8 + 9 - 10")
	h := parser.NewHistory()

	var add parser.Func
	value := combinator.Memoize("VALUE", h, combinator.Any("value",
		terminal.Integer(),
		&add,
	))

	add = combinator.Memoize("SUM", h, combinator.SepBy1(
		"SUM", h, value, combinator.Choice("+ or -", terminal.Rune('+', "+"), terminal.Rune('-', "-")),
		ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, error) {
			sum := 0
			modifier := 1
			for _, node := range nodes {
				v, err := node.Value(ctx)
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

	s := combinator.Seq(builder.Select(0), value, parser.End())
	result, err := parsley.EvaluateText([]byte(input), true, s, nil)
	require.Nil(t, err)
	assert.Equal(t, -5, result)
	assert.Equal(t, 1242, parser.Stat.GetSumCallCount())
}
