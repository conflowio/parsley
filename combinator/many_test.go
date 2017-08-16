// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator_test

import (
	"fmt"
	"testing"

	"github.com/opsidian/parsley"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/test"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Let's define a parser which accepts any number of "a" characters
func ExampleMany() {
	concat := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		var res string
		for _, node := range nodes {
			val, _ := node.Value(ctx)
			res += string(val.(rune))
		}
		return res, nil
	})
	s := combinator.Many(builder.All("a", concat), terminal.Rune('a', "A"))
	value, _ := parsley.EvaluateText([]byte("aaaaabbbbb"), true, s, nil)
	fmt.Printf("%T %v\n", value, value)
	// Output: string aaaaa
}

// Let's define a parser which accepts one or many "a" characters
func ExampleMany1() {
	concat := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		var res string
		for _, node := range nodes {
			val, _ := node.Value(ctx)
			res += string(val.(rune))
		}
		return res, nil
	})
	s := combinator.Many1(builder.All("a", concat), terminal.Rune('a', "A"))
	value, _ := parsley.EvaluateText([]byte("aaaaabbbbb"), true, s, nil)
	fmt.Printf("%T %v\n", value, value)
	// Output: string aaaaa
}

func TestManyShouldPanicIfNoBuilder(t *testing.T) {
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, nil
	})
	assert.Panics(t, func() {
		combinator.Many(nil, p)(parser.EmptyLeftRecCtx(), test.NewReader(0, 0, false, false))
	})
	assert.Panics(t, func() {
		combinator.Many1(nil, p)(parser.EmptyLeftRecCtx(), test.NewReader(0, 0, false, false))
	})
}

func TestManyShouldCombineParserResults(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	pResults := []parser.ResultSet{
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "a"), test.NewReader(1, 1, false, true)),
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "b"), test.NewReader(2, 1, false, true)),
		),
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "c"), test.NewReader(3, 1, false, true)),
		),
		nil,
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "d"), test.NewReader(4, 1, false, true)),
		),
		nil,
	}

	pi := 0
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		defer func() { pi++ }()
		if pi < len(pResults) {
			return parser.NoCurtailingParsers(), pResults[pi], nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		var res string
		for _, node := range nodes {
			val, _ := node.Value(nil)
			res += val.(string)
		}
		first := nodes[0].(ast.TerminalNode)
		return ast.NewTerminalNode("STR", first.Pos(), res)
	})

	_, rs, err := combinator.Many(nodeBuilder, p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, parser.NewResultSet(
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(3, 1, false, true)),
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "bd"), test.NewReader(4, 1, false, true)),
	), rs)
	assert.Nil(t, err)

	pi = 0 // Reset parser call index
	_, rs, err = combinator.Many1(nodeBuilder, p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, parser.NewResultSet(
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(3, 1, false, true)),
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "bd"), test.NewReader(4, 1, false, true)),
	), rs)
	assert.Nil(t, err)
}

func TestMany1ShouldReturnNoResultIfNoMatch(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "ERR1")
	})

	cp, rs, err := combinator.Many1(builder.Nil(), p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(1), err.Pos())
}

func TestManyShouldReturnEmptyResultIfNoMatch(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "ERR1")
	})

	cp, rs, err := combinator.Many(builder.Nil(), p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Equal(t, rs, parser.NewResult(nil, r).AsSet())
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(1), err.Pos())
}

func TestManyShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	pi := 0
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		defer func() { pi++ }()
		if pi == 0 {
			return data.NewIntSet(0, 1), parser.NewResult(nil, r).AsSet(), nil
		} else {
			return data.NewIntSet(1, 2), nil, reader.NewError(test.NewPosition(1), "ERR1")
		}
	})

	cp, _, _ := combinator.Many(builder.Nil(), p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)

	pi = 0 // Reset parser call index
	cp, _, _ = combinator.Many1(builder.Nil(), p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
}
