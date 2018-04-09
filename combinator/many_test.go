// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator_test

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Let's define a parser which accepts any number of "a" characters
func ExampleMany() {
	concat := ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		var res string
		for _, node := range nodes {
			val, _ := node.Value(ctx)
			res += string(val.(rune))
		}
		return res, nil
	})
	p := combinator.Many("MANY_A", terminal.Rune('a')).Bind(concat)
	r := text.NewReader(text.NewFile("example.file", []byte("aaaaa")))
	value, _ := parsley.Evaluate(parser.NewHistory(), r, combinator.Sentence(p), nil)

	fmt.Printf("%T %v\n", value, value)
	// Output: string aaaaa
}

//
// // Let's define a parser which accepts one or many "a" characters
// func ExampleMany1() {
// 	concat := parsley.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
// 		var res string
// 		for _, node := range nodes {
// 			val, _ := node.Value(ctx)
// 			res += string(val.(rune))
// 		}
// 		return res, nil
// 	})
// 	p := combinator.Many1(builder.All("a", concat), terminal.Rune('a', "A"))
// 	s := parsley.NewSentence(p)
// 	value, _, _ := s.Evaluate(text.NewReader([]byte("aaaaa"), "", true), nil)
// 	fmt.Printf("%T %v\n", value, value)
// 	// Output: string aaaaa
// }
//
// func TestManyShouldPanicIfNoBuilder(t *testing.T) {
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), nil, nil
// 	})
// 	assert.Panics(t, func() {
// 		combinator.Many(nil, p)(parser.NewHistory(), parser.EmptyLeftRecCtx(), test.NewReader(0, 0, false, false))
// 	})
// 	assert.Panics(t, func() {
// 		combinator.Many1(nil, p)(parser.NewHistory(), parser.EmptyLeftRecCtx(), test.NewReader(0, 0, false, false))
// 	})
// }
//
// func TestManyShouldCombineParserResults(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	pResults := []parser.ResultSet{
// 		parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "a"), test.NewReader(1, 1, false, true)),
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "b"), test.NewReader(2, 1, false, true)),
// 		),
// 		parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "c"), test.NewReader(3, 1, false, true)),
// 		),
// 		nil,
// 		parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "d"), test.NewReader(4, 1, false, true)),
// 		),
// 		nil,
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
// 	nodeBuilder := parsley.NodeBuilderFunc(func(nodes []parsley.Node) parsley.Node {
// 		var res string
// 		for _, node := range nodes {
// 			val, _ := node.Value(nil)
// 			res += val.(string)
// 		}
// 		first := nodes[0].(ast.TerminalNode)
// 		return ast.NewTerminalNode("STR", first.Pos(), res)
// 	})
//
// 	h := parser.NewHistory()
// 	_, rs, err := combinator.Many(nodeBuilder, p).Parse(h, parser.EmptyLeftRecCtx(), r)
// 	assert.EqualValues(t, parser.NewResultSet(
// 		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(3, 1, false, true)),
// 		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "bd"), test.NewReader(4, 1, false, true)),
// 	), rs)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 5, h.CallCount())
//
// 	pi = 0 // Reset parser call index
// 	h = parser.NewHistory()
// 	_, rs, err = combinator.Many1(nodeBuilder, p).Parse(h, parser.EmptyLeftRecCtx(), r)
// 	assert.EqualValues(t, parser.NewResultSet(
// 		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(3, 1, false, true)),
// 		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "bd"), test.NewReader(4, 1, false, true)),
// 	), rs)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 5, h.CallCount())
// }
//
// func TestMany1ShouldReturnNoResultIfNoMatch(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "ERR1")
// 	})
//
// 	cp, rs, err := combinator.Many1(builder.Nil(), p).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 	assert.Empty(t, rs)
// 	require.NotNil(t, err)
// 	assert.Equal(t, test.NewPosition(1), err.Pos())
// }
//
// func TestManyShouldReturnEmptyResultIfNoMatch(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "ERR1")
// 	})
//
// 	cp, rs, err := combinator.Many(builder.Nil(), p).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 	assert.Equal(t, rs, parser.NewResult(nil, r).AsSet())
// 	require.NotNil(t, err)
// 	assert.Equal(t, test.NewPosition(1), err.Pos())
// }
//
// func TestManyShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	pi := 0
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		defer func() { pi++ }()
// 		if pi == 0 {
// 			return data.NewIntSet(0, 1), parser.NewResult(nil, r).AsSet(), nil
// 		} else {
// 			return data.NewIntSet(1, 2), nil, reader.NewError(test.NewPosition(1), "ERR1")
// 		}
// 	})
//
// 	cp, _, _ := combinator.Many(builder.Nil(), p).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
//
// 	pi = 0 // Reset parser call index
// 	cp, _, _ = combinator.Many1(builder.Nil(), p).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
// }
