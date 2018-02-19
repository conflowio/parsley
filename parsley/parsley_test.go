// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley_test

import (
	"fmt"
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/reader/mocks"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Let's define a simple parser which is able to parse adding two integers.
func ExampleSentence_Parse() {
	add := combinator.Seq(
		builder.BinaryOperation(
			ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
				value0, _ := nodes[0].Value(ctx)
				value1, _ := nodes[1].Value(ctx)
				return value0.(int) + value1.(int), nil
			}),
		),
		terminal.Integer(),
		terminal.Rune('+', "ADD"),
		terminal.Integer(),
	)
	s := parsley.NewSentence(add)
	node, _, err := s.Parse(text.NewReader([]byte("1 + 2"), "", true))
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
func ExampleSentence_Evaluate() {
	add := combinator.Seq(
		builder.BinaryOperation(
			ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
				value0, _ := nodes[0].Value(ctx)
				value1, _ := nodes[1].Value(ctx)
				return value0.(int) + value1.(int), nil
			}),
		),
		terminal.Integer(),
		terminal.Rune('+', "ADD"),
		terminal.Integer(),
	)
	s := parsley.NewSentence(add)
	value, _, err := s.Evaluate(text.NewReader([]byte("1 + 2"), "", true), nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %d\n", value.(int))
	// Output: Result: 3
}

func TestParseShouldRunParserAndReturnNode(t *testing.T) {
	expectedNode := ast.NewTerminalNode("STRING", text.NewPosition(0, 1, 1), "RES")
	p := terminal.Word("input", "STRING", "RES")
	s := parsley.NewSentence(p)
	node, _, err := s.Parse(text.NewReader([]byte("input"), "", true))
	assert.Equal(t, expectedNode, node)
	assert.Nil(t, err)
}

func TestParseShouldHandleEmptyResult(t *testing.T) {
	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(text.NewPosition(2, 1, 3), "encountered a test error")
	})
	s := parsley.NewSentence(p)
	node, _, err := s.Parse(text.NewReader([]byte("input"), "", true))
	assert.Error(t, err)
	assert.Equal(t, "Failed to parse the input: encountered a test error at 1:3", err.Error())
	assert.Nil(t, node)
}

func TestParseShouldHandleNilNode(t *testing.T) {
	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), parser.NewResult(nil, r).AsSet(), nil
	})
	s := parsley.NewSentence(p)
	node, _, err := s.Parse(text.NewReader([]byte(""), "", true))
	assert.Nil(t, err)
	assert.Nil(t, node)
}

func TestEvaluateShouldRunParserAndReturnValue(t *testing.T) {
	expectedValue := "RES"
	p := terminal.Word("input", "STRING", "RES")
	s := parsley.NewSentence(p)
	value, _, err := s.Evaluate(text.NewReader([]byte("input"), "", true), nil)
	assert.Equal(t, expectedValue, value)
	assert.Nil(t, err)
}

func TestEvaluateShouldPassContext(t *testing.T) {
	ctx := "testctx"
	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		return ctx, nil
	})
	node := ast.NewNonTerminalNode("STRING", nil, interpreter)
	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		r.(*text.Reader).ReadMatch(".*", true)
		return parser.NoCurtailingParsers(), parser.NewResult(node, r).AsSet(), nil
	})
	s := parsley.NewSentence(p)
	value, _, err := s.Evaluate(text.NewReader([]byte("input"), "", true), ctx)
	assert.Equal(t, ctx, value)
	assert.Nil(t, err)
}

func TestEvaluateShouldHandleEmptyResult(t *testing.T) {
	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(text.NewPosition(2, 1, 3), "encountered a test error")
	})
	s := parsley.NewSentence(p)
	value, _, err := s.Evaluate(text.NewReader([]byte("input"), "", true), nil)
	assert.Error(t, err)
	assert.Equal(t, "Failed to parse the input: encountered a test error at 1:3", err.Error())
	assert.Nil(t, value)
}

func TestEvaluateShouldHandleNilNode(t *testing.T) {
	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), parser.NewResult(nil, r).AsSet(), nil
	})
	s := parsley.NewSentence(p)
	value, _, err := s.Evaluate(text.NewReader([]byte(""), "", true), nil)
	assert.Nil(t, err)
	assert.Nil(t, value)
}

func TestEvaluateShouldHandleInterpreterError(t *testing.T) {
	randomChild := ast.NewTerminalNode("X", text.NewPosition(1, 2, 3), "X")
	expectedErr := new(mocks.Error)
	node := ast.NewNonTerminalNode("ERR", []ast.Node{randomChild}, ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		return nil, expectedErr
	}))
	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		r.(*text.Reader).ReadMatch(".*", true)
		return parser.NoCurtailingParsers(), parser.NewResult(node, r).AsSet(), nil
	})
	s := parsley.NewSentence(p)
	value, _, err := s.Evaluate(text.NewReader([]byte("input"), "", true), nil)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, value)
}

func TestDirectLeftRecursion(t *testing.T) {
	input := "abbbbbbbbbbbbbbbbbbb"

	concatNodes := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
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
	a = combinator.Memoize(combinator.Any("a or ab",
		combinator.Seq(builder.All("AB", concatNodes),
			&a,
			terminal.Rune('b', "CHAR"),
		),
		terminal.Rune('a', "CHAR"),
	))
	s := parsley.NewSentence(a)
	result, h, err := s.Evaluate(text.NewReader([]byte(input), "", true), nil)
	require.Nil(t, err)
	assert.Equal(t, input, result)
	assert.Equal(t, 318, h.CallCount())
}

func TestIndirectLeftRecursion(t *testing.T) {
	input := []byte("1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10")

	var add parser.Func
	value := combinator.Memoize(combinator.Any("value",
		terminal.Integer(),
		&add,
	))

	add = combinator.Memoize(combinator.Seq(
		builder.BinaryOperation(
			ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
				value0, _ := nodes[0].Value(ctx)
				value1, _ := nodes[1].Value(ctx)
				return value0.(int) + value1.(int), nil
			}),
		),
		value,
		terminal.Rune('+', "ADD"),
		value,
	))
	s := parsley.NewSentence(add)
	result, h, err := s.Evaluate(text.NewReader([]byte(input), "", true), nil)
	require.Nil(t, err)
	assert.Equal(t, 55, result)
	assert.Equal(t, 3477, h.CallCount())
}

func TestSepBy(t *testing.T) {
	input := []byte("1 - 2 + 3 - 4 + 5 - 6 + 7 - 8 + 9 - 10")

	var add parser.Func
	value := combinator.Memoize(combinator.Any("value",
		terminal.Integer(),
		&add,
	))

	add = combinator.Memoize(combinator.SepBy1(
		"SUM", value, combinator.Choice("+ or -", terminal.Rune('+', "+"), terminal.Rune('-', "-")),
		ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
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

	s := parsley.NewSentence(value)
	result, h, err := s.Evaluate(text.NewReader([]byte(input), "", true), nil)
	require.Nil(t, err)
	assert.Equal(t, -5, result)
	assert.Equal(t, 1242, h.CallCount())
}
