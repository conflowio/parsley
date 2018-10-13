// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator_test

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Let's define a simple language where you define an integer array.
// The language would be left recursive, but using SepBy we can avoid this.
// The grammar is: S -> [I(,I)*], I -> any integer
func ExampleSepBy() {
	arr := ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		var res []int64
		for i := 0; i < len(nodes); i += 2 {
			val, _ := nodes[i].Value(ctx)
			res = append(res, val.(int64))
		}
		return res, nil
	})

	intList := combinator.SepBy("integer list", terminal.Integer(), terminal.Rune(',')).Bind(arr)
	p := combinator.Seq("ARR", "array",
		terminal.Rune('['),
		intList,
		terminal.Rune(']'),
	).Bind(interpreter.Select(1))

	r := text.NewReader(text.NewFile("example.file", []byte("[]")))
	ctx := parsley.NewContext(parsley.NewFileSet(), r)

	value1, _ := parsley.Evaluate(ctx, combinator.Sentence(p), nil)
	fmt.Printf("%T %v\n", value1, value1)

	r = text.NewReader(text.NewFile("example.file", []byte("[1,2,3]")))
	ctx = parsley.NewContext(parsley.NewFileSet(), r)
	value2, _ := parsley.Evaluate(ctx, combinator.Sentence(p), nil)
	fmt.Printf("%T %v\n", value2, value2)

	// Output: []int64 []
	// []int64 [1 2 3]
}

// Let's define a simple language where you can add integer numbers.
// The language would be left recursive, but using SepBy1 we can avoid this.
// The grammar is: S -> I(+I)*, I -> any integer
// The "<empty>" result will never be returned as the SepBy1 doesn't match zero p occurrences.
func ExampleSepBy1() {
	interpreter := ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		if len(nodes) == 0 {
			return "<empty>", nil
		}
		var sum int64
		for i := 0; i < len(nodes); i += 2 {
			val, _ := nodes[i].Value(ctx)
			sum += val.(int64)
		}
		return sum, nil
	})

	p := combinator.SepBy1("integers separated by '+'", terminal.Integer(), terminal.Rune('+')).Bind(interpreter)

	r := text.NewReader(text.NewFile("example.file", []byte("")))
	ctx := parsley.NewContext(parsley.NewFileSet(), r)
	value1, _ := parsley.Evaluate(ctx, combinator.Sentence(p), nil)
	fmt.Printf("%T %v\n", value1, value1)

	r = text.NewReader(text.NewFile("example.file", []byte("1")))
	ctx = parsley.NewContext(parsley.NewFileSet(), r)
	value2, _ := parsley.Evaluate(ctx, combinator.Sentence(p), nil)
	fmt.Printf("%T %v\n", value2, value2)

	r = text.NewReader(text.NewFile("example.file", []byte("1+2+3")))
	ctx = parsley.NewContext(parsley.NewFileSet(), r)
	value3, _ := parsley.Evaluate(ctx, combinator.Sentence(p), nil)
	fmt.Printf("%T %v\n", value3, value3)
	// Output: <nil> <nil>
	// int64 1
	// int64 6
}

//
// func TestSepByShouldCombineParserResults(t *testing.T) {
// 	pResults := []parser.ResultSet{
// 		parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(0), "a"), test.NewReader(1, 1, false, true)),
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "b"), test.NewReader(2, 1, false, true)),
// 		),
// 		parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(4), "c"), test.NewReader(5, 1, false, true)),
// 		),
// 		parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(5), "d"), test.NewReader(6, 1, false, true)),
// 		),
// 	}
//
// 	sepResults := []parser.ResultSet{
// 		parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), ","), test.NewReader(3, 1, false, true)),
// 		),
// 		nil,
// 		parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(3), ","), test.NewReader(4, 1, false, true)),
// 		),
// 	}
//
// 	pi := 0
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		defer func() { pi++ }()
// 		if pi < len(pResults) {
// 			return parser.NoCurtailingParsers(), pResults[pi], nil
// 		} else {
// 			return parser.NoCurtailingParsers(), nil, nil
// 		}
// 	})
//
// 	sepi := 0
// 	sep := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		defer func() { sepi++ }()
// 		if sepi < len(sepResults) {
// 			return parser.NoCurtailingParsers(), sepResults[sepi], nil
// 		} else {
// 			return parser.NoCurtailingParsers(), nil, nil
// 		}
// 	})
//
// 	interpreter := parsley.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
// 		res := ""
// 		for _, node := range nodes {
// 			value, _ := node.Value(ctx)
// 			res += value.(string) + "|"
// 		}
// 		return res, nil
// 	})
//
// 	parsers := []parsley.Parser{
// 		combinator.SepBy("TEST", p, sep, interpreter),
// 		combinator.SepBy1("TEST", p, sep, interpreter),
// 		combinator.SepByOrValue("TEST", p, sep, interpreter),
// 		combinator.SepByOrValue1("TEST", p, sep, interpreter),
// 	}
//
// 	for _, p := range parsers {
// 		pi = 0
// 		sepi = 0
// 		r := test.NewReader(0, 1, false, false)
// 		h := parser.NewHistory()
// 		_, rs, err := p.Parse(h, parser.EmptyLeftRecCtx(), r)
// 		require.Len(t, rs, 2)
// 		val0, _ := rs[0].Node().Value(nil)
// 		val1, _ := rs[1].Node().Value(nil)
// 		assert.Equal(t, "a|,|c|", val0)
// 		assert.Equal(t, "b|,|d|", val1)
// 		assert.Nil(t, err)
// 		assert.Equal(t, 13, h.CallCount())
// 	}
// }
//
// func TestSepByShouldNotFlattenNonTerminals(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	interpreter1 := parsley.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
// 		res := ""
// 		for _, node := range nodes {
// 			value, _ := node.Value(ctx)
// 			res += value.(string) + "&"
// 		}
// 		return res, nil
// 	})
//
// 	pResults := []parser.ResultSet{
// 		parser.NewResultSet(
// 			parser.NewResult(
// 				ast.NewNonTerminalNode("NT", []parsley.Node{
// 					ast.NewTerminalNode("STR", test.NewPosition(0), "a"),
// 					ast.NewTerminalNode("STR", test.NewPosition(1), "b"),
// 				}, interpreter1),
// 				test.NewReader(2, 1, false, true),
// 			),
// 		),
// 		parser.NewResultSet(
// 			parser.NewResult(
// 				ast.NewNonTerminalNode("NT", []parsley.Node{
// 					ast.NewTerminalNode("STR", test.NewPosition(3), "c"),
// 					ast.NewTerminalNode("STR", test.NewPosition(4), "d"),
// 				}, interpreter1),
// 				test.NewReader(5, 1, false, true),
// 			),
// 		),
// 	}
//
// 	sepResults := []parser.ResultSet{
// 		parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), ","), test.NewReader(3, 1, false, true)),
// 		),
// 	}
//
// 	pi := 0
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		defer func() { pi++ }()
// 		if pi < len(pResults) {
// 			return parser.NoCurtailingParsers(), pResults[pi], nil
// 		} else {
// 			return parser.NoCurtailingParsers(), nil, nil
// 		}
// 	})
//
// 	sepi := 0
// 	sep := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		defer func() { sepi++ }()
// 		if sepi < len(sepResults) {
// 			return parser.NoCurtailingParsers(), sepResults[sepi], nil
// 		} else {
// 			return parser.NoCurtailingParsers(), nil, nil
// 		}
// 	})
//
// 	interpreter2 := parsley.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
// 		res := ""
// 		for _, node := range nodes {
// 			value, _ := node.Value(ctx)
// 			res += value.(string) + "|"
// 		}
// 		return res, nil
// 	})
//
// 	parsers := []parsley.Parser{
// 		combinator.SepBy("TEST", p, sep, interpreter2),
// 		combinator.SepBy1("TEST", p, sep, interpreter2),
// 		combinator.SepByOrValue("TEST", p, sep, interpreter2),
// 		combinator.SepByOrValue1("TEST", p, sep, interpreter2),
// 	}
//
// 	for _, p := range parsers {
// 		pi = 0
// 		sepi = 0
// 		_, rs, _ := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 		require.Len(t, rs, 1)
// 		val0, _ := rs[0].Node().Value(nil)
// 		assert.Equal(t, "a&b&|,|c&d&|", val0)
// 	}
// }
//
// func TestSepByShouldReturnEmptyResultIfNoMatch(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "TEST1")
// 	})
//
// 	parsers := []parsley.Parser{
// 		combinator.SepBy("TEST", p, p, nil),
// 		combinator.SepByOrValue("TEST", p, p, nil),
// 	}
//
// 	for _, p := range parsers {
// 		cp, rs, err := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 		assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 		assert.Equal(t, parser.NewResult(ast.NewNonTerminalNode("TEST", nil, nil), r).AsSet(), rs)
// 		require.NotNil(t, err)
// 		assert.Equal(t, test.NewPosition(1), err.Pos())
// 	}
// }
//
// func TestSepBy1ShouldReturnNilIfNoResult(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "TEST1")
// 	})
//
// 	parsers := []parsley.Parser{
// 		combinator.SepBy1("TEST", p, p, nil),
// 		combinator.SepByOrValue1("TEST", p, p, nil),
// 	}
//
// 	for _, p := range parsers {
// 		cp, rs, err := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 		assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 		assert.Empty(t, rs)
// 		require.NotNil(t, err)
// 		assert.Equal(t, test.NewPosition(1), err.Pos())
// 	}
// }
//
// func TestSepByShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	pi := 0
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		defer func() { pi++ }()
// 		if pi < 1 {
// 			return data.NewIntSet(0, 1), parser.NewResult(nil, r).AsSet(), nil
// 		} else {
// 			return parser.NoCurtailingParsers(), nil, nil
// 		}
// 	})
//
// 	sepi := 0
// 	sep := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		defer func() { sepi++ }()
// 		if sepi < 1 {
// 			return data.NewIntSet(1, 2), nil, nil
// 		} else {
// 			return parser.NoCurtailingParsers(), nil, nil
// 		}
// 	})
//
// 	parsers := []parsley.Parser{
// 		combinator.SepBy("TEST", p, sep, nil),
// 		combinator.SepBy1("TEST", p, sep, nil),
// 		combinator.SepByOrValue("TEST", p, sep, nil),
// 		combinator.SepByOrValue1("TEST", p, sep, nil),
// 	}
//
// 	for _, p := range parsers {
// 		pi = 0
// 		sepi = 0
// 		cp, _, _ := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 		assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
// 	}
// }
//
// func TestSepByShouldNotReturnValueIfOnlyOne(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
// 	value := ast.NewTerminalNode("VALUE", test.NewPosition(0), "value")
//
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), parser.NewResult(value, r).AsSet(), nil
// 	})
//
// 	s := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), nil, nil
// 	})
//
// 	expectedNode := ast.NewNonTerminalNode("TEST", []parsley.Node{value}, nil)
//
// 	parsers := []parsley.Parser{
// 		combinator.SepBy("TEST", p, s, nil),
// 		combinator.SepBy1("TEST", p, s, nil),
// 	}
//
// 	for _, p := range parsers {
// 		cp, rs, err := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 		assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 		assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), rs)
// 		assert.Nil(t, err)
// 	}
// }
//
// func TestSepByOrValueShouldReturnValueIfOnlyOne(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
// 	value := ast.NewTerminalNode("VALUE", test.NewPosition(0), "value")
//
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), parser.NewResult(value, r).AsSet(), nil
// 	})
//
// 	s := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), nil, nil
// 	})
//
// 	parsers := []parsley.Parser{
// 		combinator.SepByOrValue("TEST", p, s, nil),
// 		combinator.SepByOrValue1("TEST", p, s, nil),
// 	}
//
// 	for _, p := range parsers {
// 		cp, rs, err := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 		assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 		assert.Equal(t, parser.NewResult(value, r).AsSet(), rs)
// 		assert.Nil(t, err)
// 	}
// }
