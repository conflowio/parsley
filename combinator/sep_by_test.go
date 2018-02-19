// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator_test

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
	"github.com/opsidian/parsley/test"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Let's define a simple language where you define an integer array.
// The language would be left recursive, but using SepBy (which is using Many and Seq) we can avoid this.
// The grammar is: S -> [I(,I)*], I -> any integer
func ExampleSepBy() {
	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		var res []int
		for i := 0; i < len(nodes); i += 2 {
			val, _ := nodes[i].Value(ctx)
			res = append(res, val.(int))
		}
		return res, nil
	})

	intList := combinator.SepBy("ARR", terminal.Integer(), terminal.Rune(',', "SEP"), interpreter)
	p := combinator.Seq(builder.Select(1), terminal.Rune('[', "ARR_START"), intList, terminal.Rune(']', "ARR_END"))
	s := parsley.NewSentence(p)

	value1, _, _ := s.Evaluate(text.NewReader([]byte("[]"), "", true), nil)
	fmt.Printf("%T %v\n", value1, value1)

	value2, _, _ := s.Evaluate(text.NewReader([]byte("[1, 2, 3]"), "", true), nil)
	fmt.Printf("%T %v\n", value2, value2)
	// Output: []int []
	// []int [1 2 3]
}

// Let's define a simple language where you can add integer numbers.
// The language would be left recursive, but using SepBy1 (which is using Many and Seq) we can avoid this.
// The grammar is: S -> I(+I)*, I -> any integer
func ExampleSepBy1() {
	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		sum := 0
		for i := 0; i < len(nodes); i += 2 {
			val, _ := nodes[i].Value(ctx)
			sum += val.(int)
		}
		return sum, nil
	})

	p := combinator.SepBy1("SUM", terminal.Integer(), terminal.Rune('+', "+"), interpreter)
	s := parsley.NewSentence(p)
	value1, _, _ := s.Evaluate(text.NewReader([]byte("1"), "", true), nil)
	fmt.Printf("%T %v\n", value1, value1)

	value2, _, _ := s.Evaluate(text.NewReader([]byte("1 + 2 + 3"), "", true), nil)
	fmt.Printf("%T %v\n", value2, value2)
	// Output: int 1
	// int 6
}

func TestSepByShouldCombineParserResults(t *testing.T) {
	pResults := []parser.ResultSet{
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(0), "a"), test.NewReader(1, 1, false, true)),
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "b"), test.NewReader(2, 1, false, true)),
		),
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(4), "c"), test.NewReader(5, 1, false, true)),
		),
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(5), "d"), test.NewReader(6, 1, false, true)),
		),
	}

	sepResults := []parser.ResultSet{
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), ","), test.NewReader(3, 1, false, true)),
		),
		nil,
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(3), ","), test.NewReader(4, 1, false, true)),
		),
	}

	pi := 0
	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		defer func() { pi++ }()
		if pi < len(pResults) {
			return parser.NoCurtailingParsers(), pResults[pi], nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	sepi := 0
	sep := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		defer func() { sepi++ }()
		if sepi < len(sepResults) {
			return parser.NoCurtailingParsers(), sepResults[sepi], nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		res := ""
		for _, node := range nodes {
			value, _ := node.Value(ctx)
			res += value.(string) + "|"
		}
		return res, nil
	})

	parsers := []parser.Parser{
		combinator.SepBy("TEST", p, sep, interpreter),
		combinator.SepBy1("TEST", p, sep, interpreter),
		combinator.SepByOrValue("TEST", p, sep, interpreter),
		combinator.SepByOrValue1("TEST", p, sep, interpreter),
	}

	for _, p := range parsers {
		pi = 0
		sepi = 0
		r := test.NewReader(0, 1, false, false)
		h := parser.NewHistory()
		_, rs, err := p.Parse(h, parser.EmptyLeftRecCtx(), r)
		require.Len(t, rs, 2)
		val0, _ := rs[0].Node().Value(nil)
		val1, _ := rs[1].Node().Value(nil)
		assert.Equal(t, "a|,|c|", val0)
		assert.Equal(t, "b|,|d|", val1)
		assert.Nil(t, err)
		assert.Equal(t, 13, h.CallCount())
	}
}

func TestSepByShouldNotFlattenNonTerminals(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	interpreter1 := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		res := ""
		for _, node := range nodes {
			value, _ := node.Value(ctx)
			res += value.(string) + "&"
		}
		return res, nil
	})

	pResults := []parser.ResultSet{
		parser.NewResultSet(
			parser.NewResult(
				ast.NewNonTerminalNode("NT", []ast.Node{
					ast.NewTerminalNode("STR", test.NewPosition(0), "a"),
					ast.NewTerminalNode("STR", test.NewPosition(1), "b"),
				}, interpreter1),
				test.NewReader(2, 1, false, true),
			),
		),
		parser.NewResultSet(
			parser.NewResult(
				ast.NewNonTerminalNode("NT", []ast.Node{
					ast.NewTerminalNode("STR", test.NewPosition(3), "c"),
					ast.NewTerminalNode("STR", test.NewPosition(4), "d"),
				}, interpreter1),
				test.NewReader(5, 1, false, true),
			),
		),
	}

	sepResults := []parser.ResultSet{
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), ","), test.NewReader(3, 1, false, true)),
		),
	}

	pi := 0
	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		defer func() { pi++ }()
		if pi < len(pResults) {
			return parser.NoCurtailingParsers(), pResults[pi], nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	sepi := 0
	sep := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		defer func() { sepi++ }()
		if sepi < len(sepResults) {
			return parser.NoCurtailingParsers(), sepResults[sepi], nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	interpreter2 := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		res := ""
		for _, node := range nodes {
			value, _ := node.Value(ctx)
			res += value.(string) + "|"
		}
		return res, nil
	})

	parsers := []parser.Parser{
		combinator.SepBy("TEST", p, sep, interpreter2),
		combinator.SepBy1("TEST", p, sep, interpreter2),
		combinator.SepByOrValue("TEST", p, sep, interpreter2),
		combinator.SepByOrValue1("TEST", p, sep, interpreter2),
	}

	for _, p := range parsers {
		pi = 0
		sepi = 0
		_, rs, _ := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
		require.Len(t, rs, 1)
		val0, _ := rs[0].Node().Value(nil)
		assert.Equal(t, "a&b&|,|c&d&|", val0)
	}
}

func TestSepByShouldReturnEmptyResultIfNoMatch(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "TEST1")
	})

	parsers := []parser.Parser{
		combinator.SepBy("TEST", p, p, nil),
		combinator.SepByOrValue("TEST", p, p, nil),
	}

	for _, p := range parsers {
		cp, rs, err := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
		assert.Equal(t, parser.NoCurtailingParsers(), cp)
		assert.Equal(t, parser.NewResult(ast.NewNonTerminalNode("TEST", nil, nil), r).AsSet(), rs)
		require.NotNil(t, err)
		assert.Equal(t, test.NewPosition(1), err.Pos())
	}
}

func TestSepBy1ShouldReturnNilIfNoResult(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "TEST1")
	})

	parsers := []parser.Parser{
		combinator.SepBy1("TEST", p, p, nil),
		combinator.SepByOrValue1("TEST", p, p, nil),
	}

	for _, p := range parsers {
		cp, rs, err := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
		assert.Equal(t, parser.NoCurtailingParsers(), cp)
		assert.Empty(t, rs)
		require.NotNil(t, err)
		assert.Equal(t, test.NewPosition(1), err.Pos())
	}
}

func TestSepByShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	pi := 0
	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		defer func() { pi++ }()
		if pi < 1 {
			return data.NewIntSet(0, 1), parser.NewResult(nil, r).AsSet(), nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	sepi := 0
	sep := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		defer func() { sepi++ }()
		if sepi < 1 {
			return data.NewIntSet(1, 2), nil, nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	parsers := []parser.Parser{
		combinator.SepBy("TEST", p, sep, nil),
		combinator.SepBy1("TEST", p, sep, nil),
		combinator.SepByOrValue("TEST", p, sep, nil),
		combinator.SepByOrValue1("TEST", p, sep, nil),
	}

	for _, p := range parsers {
		pi = 0
		sepi = 0
		cp, _, _ := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
		assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
	}
}

func TestSepByShouldNotReturnValueIfOnlyOne(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	value := ast.NewTerminalNode("VALUE", test.NewPosition(0), "value")

	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), parser.NewResult(value, r).AsSet(), nil
	})

	s := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, nil
	})

	expectedNode := ast.NewNonTerminalNode("TEST", []ast.Node{value}, nil)

	parsers := []parser.Parser{
		combinator.SepBy("TEST", p, s, nil),
		combinator.SepBy1("TEST", p, s, nil),
	}

	for _, p := range parsers {
		cp, rs, err := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
		assert.Equal(t, parser.NoCurtailingParsers(), cp)
		assert.Equal(t, parser.NewResult(expectedNode, r.Clone()).AsSet(), rs)
		assert.Nil(t, err)
	}
}

func TestSepByOrValueShouldReturnValueIfOnlyOne(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	value := ast.NewTerminalNode("VALUE", test.NewPosition(0), "value")

	p := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), parser.NewResult(value, r).AsSet(), nil
	})

	s := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, nil
	})

	parsers := []parser.Parser{
		combinator.SepByOrValue("TEST", p, s, nil),
		combinator.SepByOrValue1("TEST", p, s, nil),
	}

	for _, p := range parsers {
		cp, rs, err := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
		assert.Equal(t, parser.NoCurtailingParsers(), cp)
		assert.Equal(t, parser.NewResult(value, r.Clone()).AsSet(), rs)
		assert.Nil(t, err)
	}
}
