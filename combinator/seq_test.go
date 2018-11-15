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
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Let's define a parser which accepts "a", "b", "c" characters in order.
func ExampleSeqOf() {
	concat := ast.InterpreterFunc(func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
		var res string
		for _, node := range node.Children() {
			val, _ := node.Value(userCtx)
			res += string(val.(rune))
		}
		return res, nil
	})

	p := combinator.SeqOf(
		terminal.Rune('a'),
		terminal.Rune('b'),
		terminal.Rune('c'),
	).Bind(concat)
	r := text.NewReader(text.NewFile("example.file", []byte("abc")))
	ctx := parsley.NewContext(parsley.NewFileSet(), r)
	value, _ := parsley.Evaluate(ctx, combinator.Sentence(p))
	fmt.Printf("%T %v\n", value, value)
	// Output: string abc
}

// Let's define a parser which accepts any prefix of the "abc" string.
func ExampleSeqTry() {
	concat := ast.InterpreterFunc(func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
		var res string
		for _, node := range node.Children() {
			val, _ := node.Value(userCtx)
			res += string(val.(rune))
		}
		return res, nil
	})

	p := combinator.SeqTry(
		terminal.Rune('a'),
		terminal.Rune('b'),
		terminal.Rune('c'),
	).Bind(concat)
	r := text.NewReader(text.NewFile("example.file", []byte("ab")))
	ctx := parsley.NewContext(parsley.NewFileSet(), r)
	value, _ := parsley.Evaluate(ctx, combinator.Sentence(p))
	fmt.Printf("%T %v\n", value, value)
	// Output: string ab
}

// Let's define a parser which does a simple integer addition
func ExampleSeqFirstOrAll() {
	add := ast.InterpreterFunc(func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
		val1, _ := node.Children()[0].Value(userCtx)
		val2, _ := node.Children()[2].Value(userCtx)
		return val1.(int64) + val2.(int64), nil
	})

	p := combinator.SeqFirstOrAll(
		terminal.Integer(),
		terminal.Rune('+'),
		terminal.Integer(),
	).Bind(add)

	r := text.NewReader(text.NewFile("example.file", []byte("1")))
	ctx := parsley.NewContext(parsley.NewFileSet(), r)
	value, _ := parsley.Evaluate(ctx, combinator.Sentence(p))
	fmt.Printf("%T %v\n", value, value)

	r2 := text.NewReader(text.NewFile("example.file", []byte("1+2")))
	ctx2 := parsley.NewContext(parsley.NewFileSet(), r2)
	value2, _ := parsley.Evaluate(ctx2, combinator.Sentence(p))
	fmt.Printf("%T %v\n", value2, value2)
	// Output: int64 1
	// int64 3
}

//
// func TestSeqShouldPanicIfNoBuilder(t *testing.T) {
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), nil, nil
// 	})
// 	assert.Panics(t, func() {
// 		combinator.Seq(nil, p)(parser.NewHistory(), parser.EmptyLeftRecCtx(), test.NewReader(0, 0, false, false))
// 	})
// 	assert.Panics(t, func() {
// 		combinator.SeqTry(nil, 0, p)(parser.NewHistory(), parser.EmptyLeftRecCtx(), test.NewReader(0, 0, false, false))
// 	})
// }
//
// func TestSeqShouldPanicIfNoParserWasGiven(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
// 	ctx := parser.EmptyLeftRecCtx()
// 	assert.Panics(t, func() { combinator.Seq(nil).Parse(parser.NewHistory(), ctx, r) })
// 	assert.Panics(t, func() { combinator.SeqTry(nil, 1).Parse(parser.NewHistory(), ctx, r) })
// }
//
// func TestSeqShouldHandleOnlyOneParser(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
// 	expectedRS := parser.NewResult(ast.NewTerminalNode("CHAR", r.Cursor(), 'a'), r).AsSet()
// 	expectedCP := data.NewIntSet(1)
// 	p1 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return expectedCP, expectedRS, nil
// 	})
//
// 	h := parser.NewHistory()
// 	cp, rs, err := combinator.Seq(builder.Select(0), p1).Parse(h, parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, expectedCP, cp)
// 	assert.Equal(t, expectedRS, rs)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 1, h.CallCount())
//
// 	h = parser.NewHistory()
// 	cp, rs, err = combinator.SeqTry(builder.Select(0), 0, p1).Parse(h, parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, expectedCP, cp)
// 	assert.Equal(t, expectedRS, rs)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 1, h.CallCount())
// }
//
// func TestSeqShouldCombineParserResults(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	p1 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "a"), test.NewReader(1, 1, false, true)),
// 			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ab"), test.NewReader(2, 1, false, true)),
// 		), nil
// 	})
// 	p2First := true
// 	p2 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		if p2First {
// 			p2First = false
// 			return parser.NoCurtailingParsers(), parser.NewResultSet(
// 				parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), "b"), test.NewReader(3, 1, false, true)),
// 				parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), "c"), test.NewReader(4, 1, false, true)),
// 			), nil
// 		} else {
// 			return parser.NoCurtailingParsers(), parser.NewResultSet(
// 				parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), "d"), test.NewReader(5, 1, false, true)),
// 			), nil
// 		}
// 	})
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
// 	_, rs, err := combinator.Seq(nodeBuilder, p1, p2).Parse(h, parser.EmptyLeftRecCtx(), r)
// 	assert.EqualValues(t, parser.NewResultSet(
// 		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ab"), test.NewReader(3, 1, false, true)),
// 		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(4, 1, false, true)),
// 		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "abd"), test.NewReader(5, 1, false, true)),
// 	), rs)
// 	assert.Nil(t, err)
//
// 	assert.EqualValues(t, 3, h.CallCount())
//
// 	h = parser.NewHistory()
// 	p2First = true
// 	_, rs, err = combinator.SeqTry(nodeBuilder, 0, p1, p2).Parse(h, parser.EmptyLeftRecCtx(), r)
// 	assert.EqualValues(t, parser.NewResultSet(
// 		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ab"), test.NewReader(3, 1, false, true)),
// 		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(4, 1, false, true)),
// 		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "abd"), test.NewReader(5, 1, false, true)),
// 	), rs)
// 	assert.Nil(t, err)
//
// 	assert.EqualValues(t, 3, h.CallCount())
// }
//
// func TestSeqShouldHandleNilResults(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	p1 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x'), r).AsSet(), nil
// 	})
//
// 	p2 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "ERR1")
// 	})
//
// 	cp, rs, err := combinator.Seq(builder.Nil(), p1, p2).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 	assert.Empty(t, rs)
// 	require.NotNil(t, err)
// 	assert.Equal(t, test.NewPosition(1), err.Pos())
//
// 	cp, rs, err = combinator.SeqTry(builder.Nil(), 2, p1, p2).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 	assert.Empty(t, rs)
// 	require.NotNil(t, err)
// 	assert.Equal(t, test.NewPosition(1), err.Pos())
// }
//
// func TestSeqTryShouldMatchLongestSequence(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	res := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x'), r)
//
// 	p1 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), res.AsSet(), nil
// 	})
//
// 	p2 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "ERR1")
// 	})
//
// 	cp, rs, err := combinator.SeqTry(builder.All("TEST", nil), 1, p1, p2).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 	assert.Equal(t, parser.NewResult(ast.NewNonTerminalNode("TEST", []parsley.Node{ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x')}, nil), r).AsSet(), rs)
// 	require.NotNil(t, err)
// 	assert.Equal(t, test.NewPosition(1), err.Pos())
// }
//
// func TestSeqShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	p1 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return data.NewIntSet(0, 1), parser.NewResult(nil, r).AsSet(), nil
// 	})
//
// 	p2 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return data.NewIntSet(1, 2), nil, reader.NewError(test.NewPosition(1), "ERR1")
// 	})
//
// 	cp, _, _ := combinator.Seq(builder.Nil(), p1, p2).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
//
// 	cp, _, _ = combinator.SeqTry(builder.Nil(), 0, p1, p2).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
// }
//
// func TestSeqShouldStopIfEOFTokenReached(t *testing.T) {
// 	r := test.NewReader(0, 2, false, false)
//
// 	p1 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), r).AsSet(), nil
// 	})
//
// 	p2 := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return parser.NoCurtailingParsers(), parser.NewResultSet(
// 			parser.NewResult(ast.NewTerminalNode(ast.EOF, test.NewPosition(2), nil), test.NewReader(2, 0, false, true)),
// 			parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'b'), test.NewReader(1, 1, false, true)),
// 		), nil
// 	})
//
// 	_, rs, _ := combinator.Seq(builder.Select(0), p1, p2).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.EqualValues(t,
// 		parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), test.NewReader(2, 0, false, true)).AsSet(),
// 		rs,
// 	)
//
// 	_, rs, _ = combinator.SeqTry(builder.Select(0), 0, p1, p2).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.EqualValues(t,
// 		parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), test.NewReader(2, 0, false, true)).AsSet(),
// 		rs,
// 	)
// }
