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

// Let's define a left-recursive language where we need to curtail left-recursion
// and also cache previous parser matches with Memoize.
// Grammar: S -> A, A -> a | Ab
func ExampleMemoize() {
	concat := ast.InterpreterFunc(func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
		var res string
		for _, node := range node.Children() {
			val, _ := node.Value(userCtx)
			if runeVal, ok := val.(rune); ok {
				res += string(runeVal)
			} else {
				res += val.(string)
			}
		}
		return res, nil
	})

	var p parser.Func
	p = combinator.Memoize(combinator.Any(
		terminal.Rune('a'),
		combinator.SeqOf(&p, terminal.Rune('b')).Bind(concat),
	))
	f := text.NewFile("example.file", []byte("abbbbbbbb"))
	r := text.NewReader(f)
	s := combinator.Sentence(&p)
	ctx := parsley.NewContext(parsley.NewFileSet(), r)
	value, _ := parsley.Evaluate(ctx, s)
	fmt.Printf("%T %v\n", value, value)
	// Output: string abbbbbbbb
}

//
// func TestRegisterResultShouldSaveResultForPosition(t *testing.T) {
// 	h := parser.NewHistory()
// 	parserIndex := 1
// 	node := ast.NewTerminalNode("t", test.NewPosition(0), nil)
// 	cp := parser.NoCurtailingParsers()
// 	rs := parser.NewResult(node, nil).AsSet()
// 	err := reader.NewError(test.NewPosition(1), "ERR1")
// 	h.RegisterResults(parserIndex, 2, cp, rs, err, parser.EmptyLeftRecCtx())
//
// 	actualCP, actualRS, actualErr, ok := h.GetResults(parserIndex, 2, parser.EmptyLeftRecCtx())
// 	assert.Equal(t, cp, actualCP)
// 	assert.Equal(t, rs, actualRS)
// 	assert.Equal(t, err, actualErr)
// 	assert.True(t, ok)
// }
//
// func TestRegisterResultShouldReturnNilResult(t *testing.T) {
// 	h := parser.NewHistory()
// 	parserIndex := 1
// 	err := reader.NewError(test.NewPosition(1), "ERR1")
// 	h.RegisterResults(parserIndex, 2, parser.NoCurtailingParsers(), nil, err, parser.EmptyLeftRecCtx())
// 	cp, rs, actualErr, ok := h.GetResults(parserIndex, 2, parser.EmptyLeftRecCtx())
// 	assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 	assert.Nil(t, rs)
// 	assert.Equal(t, err, actualErr)
// 	assert.True(t, ok)
// }
//
// func TestRegisterResultShouldReturnFalseWhenNoResultWasRegistered(t *testing.T) {
// 	h := parser.NewHistory()
// 	parserIndex := 1
// 	cp, rs, err, ok := h.GetResults(parserIndex, 2, parser.EmptyLeftRecCtx())
// 	assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 	assert.Nil(t, rs)
// 	assert.Nil(t, err)
// 	assert.False(t, ok)
// }
//
// func TestRegisterResultShouldHandleMultipleParsers(t *testing.T) {
// 	h := parser.NewHistory()
// 	p1Index := 1
// 	p2Index := 2
// 	node := ast.NewTerminalNode("t", test.NewPosition(0), nil)
// 	cp1 := parser.NoCurtailingParsers()
// 	cp2 := data.NewIntSet(1)
// 	rs1 := parser.NewResult(node, nil).AsSet()
// 	var rs2 parser.ResultSet
// 	err1 := reader.NewError(test.NewPosition(1), "ERR1")
// 	var err2 parsley.Error
// 	h.RegisterResults(p1Index, 1, cp1, rs1, err1, parser.EmptyLeftRecCtx())
// 	h.RegisterResults(p2Index, 2, cp2, rs2, err2, parser.EmptyLeftRecCtx())
//
// 	actualCP, actualRS, actualErr, ok := h.GetResults(p1Index, 1, parser.EmptyLeftRecCtx())
// 	assert.Equal(t, cp1, actualCP)
// 	assert.Equal(t, rs1, actualRS)
// 	assert.Equal(t, err1, actualErr)
// 	assert.True(t, ok)
//
// 	actualCP, actualRS, actualErr, ok = h.GetResults(p2Index, 2, parser.EmptyLeftRecCtx())
// 	assert.Equal(t, cp2, actualCP)
// 	assert.Equal(t, rs2, actualRS)
// 	assert.Equal(t, err2, actualErr)
// 	assert.True(t, ok)
// }
//
// func TestGetResultsShouldNotReturnCurtailedResult(t *testing.T) {
// 	h := parser.NewHistory()
// 	p1Index := 1
// 	p2Index := 2
// 	ctx := data.NewIntMap(map[int]int{
// 		p1Index: 2,
// 		p2Index: 1,
// 	})
// 	cp := data.NewIntSet(p1Index)
// 	h.RegisterResults(p1Index, 1, cp, nil, nil, ctx)
//
// 	ctx = data.NewIntMap(map[int]int{
// 		p1Index: 1,
// 		p2Index: 1,
// 	})
// 	cp, rs, err, found := h.GetResults(p1Index, 1, ctx)
// 	assert.Equal(t, parser.NoCurtailingParsers(), cp)
// 	assert.Nil(t, rs)
// 	assert.Nil(t, err)
// 	assert.False(t, found)
// }
//
// func TestGetResultsShouldReturnCurtailedResult(t *testing.T) {
// 	h := parser.NewHistory()
// 	p1Index := 1
// 	p2Index := 2
// 	ctx := data.NewIntMap(map[int]int{
// 		p1Index: 2,
// 		p2Index: 1,
// 	})
// 	cp := data.NewIntSet(p1Index)
// 	rs := parser.NewResult(nil, nil).AsSet()
// 	h.RegisterResults(p1Index, 1, cp, rs, nil, ctx)
//
// 	ctx = data.NewIntMap(map[int]int{
// 		p1Index: 1,
// 		p2Index: 1,
// 	})
//
// 	ctx = ctx.Inc(p1Index)
// 	actualCP, actualRS, err, found := h.GetResults(p1Index, 1, ctx)
// 	assert.Equal(t, cp, actualCP)
// 	assert.Equal(t, rs, actualRS)
// 	assert.Nil(t, err)
// 	assert.True(t, found)
// }
//
// func TestMemoizeShouldIncreaseLeftRecCtx(t *testing.T) {
// 	r := test.NewReader(0, 2, false, false)
// 	leftRecCtx := parser.EmptyLeftRecCtx()
// 	parserIndex := 1
// 	assert.Equal(t, leftRecCtx.Get(parserIndex), 0)
//
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		assert.Equal(t, leftRecCtx.Get(parserIndex), 1)
// 		return parser.NoCurtailingParsers(), nil, nil
// 	})
// 	combinator.Memoize(p).Parse(parser.NewHistory(), leftRecCtx, r)
// }
//
// func TestMemoizeShouldReturnParserResult(t *testing.T) {
// 	r := test.NewReader(0, 2, false, false)
//
// 	node := ast.NewTerminalNode("a", test.NewPosition(1), "a")
// 	expectedCP := data.NewIntSet(1)
// 	expectedRS := parser.NewResult(node, r).AsSet()
// 	expectedErr := reader.NewError(test.NewPosition(1), "ERR1")
//
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return expectedCP, expectedRS, expectedErr
// 	})
// 	cp, rs, err := combinator.Memoize(p).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, expectedCP, cp)
// 	assert.Equal(t, expectedRS, rs)
// 	assert.Equal(t, expectedErr, err)
// }
//
// func TestMemoizeShouldRememberResult(t *testing.T) {
// 	r := test.NewReader(0, 2, false, false)
// 	h := parser.NewHistory()
// 	node := ast.NewTerminalNode("a", test.NewPosition(1), "a")
// 	expectedCP := data.NewIntSet(1)
// 	expectedRS := parser.NewResult(node, r).AsSet()
// 	expectedErr := reader.NewError(test.NewPosition(1), "ERR1")
//
// 	called := false
// 	p := combinator.Memoize(parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		called = true
// 		return expectedCP, expectedRS, expectedErr
// 	}))
//
// 	p.Parse(h, parser.EmptyLeftRecCtx(), r)
//
// 	called = false
//
// 	cp, rs, err := p.Parse(h, parser.EmptyLeftRecCtx(), r)
//
// 	assert.False(t, called, "The parser function should not have been called")
// 	assert.Equal(t, expectedCP, cp)
// 	assert.Equal(t, expectedRS, rs)
// 	assert.Equal(t, expectedErr, err)
//
// 	p.Parse(h, parser.EmptyLeftRecCtx(), test.NewReader(1, 1, false, true))
// 	assert.True(t, called, "The parser function should have been called for the new position")
// }
//
// // The parser calls should be curtailed if the parser was already called more times than the remaining token count plus one
// func TestMemoizeShouldCurtailResult(t *testing.T) {
// 	r := test.NewReader(0, 1, false, false)
//
// 	callCount := 0
// 	var p parser.Func
// 	p = combinator.Memoize(parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		callCount++
// 		return (&p).Parse(h, leftRecCtx, r)
// 	}))
// 	cp, rs, err := p.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, 3, callCount)
// 	assert.Equal(t, 1, cp.Len())
// 	assert.Nil(t, rs)
// 	assert.Nil(t, err)
// }
