// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

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
)

// Let's define a left-recursive language where we need to curtail left-recursion
// and also cache previous parser matches with Memoize.
// Grammar: S -> A, A -> a | Ab
func ExampleHistory_Memoize() {
	h := parser.NewHistory()

	concat := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		var res string
		for _, node := range nodes {
			val, _ := node.Value(ctx)
			if runeVal, ok := val.(rune); ok {
				res += string(runeVal)
			} else {
				res += val.(string)
			}
		}
		return res, nil
	})

	var a parser.Func
	a = h.Memoize(combinator.Any("a or ab",
		terminal.Rune('a', "CHAR"),
		combinator.Seq(builder.All("AB", concat),
			&a,
			terminal.Rune('b', "CHAR"),
		),
	))
	s := combinator.Seq(builder.Select(0), &a, parser.End())

	value, _ := parsley.EvaluateText([]byte("abbbbbbbb"), true, s, nil)
	fmt.Printf("%T %v\n", value, value)
	// Output: string abbbbbbbb
}

func TestRegisterResultShouldSaveResultForPosition(t *testing.T) {
	h := parser.NewHistory()
	parserIndex := 1
	node := ast.NewTerminalNode("t", test.NewPosition(0), nil)
	cp := parser.NoCurtailingParsers()
	rs := parser.NewResult(node, nil).AsSet()
	err := reader.NewError(test.NewPosition(1), "ERR1")
	h.RegisterResults(parserIndex, 2, cp, rs, err, parser.EmptyLeftRecCtx())

	actualCP, actualRS, actualErr, ok := h.GetResults(parserIndex, 2, parser.EmptyLeftRecCtx())
	assert.Equal(t, cp, actualCP)
	assert.Equal(t, rs, actualRS)
	assert.Equal(t, err, actualErr)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnNilResult(t *testing.T) {
	h := parser.NewHistory()
	parserIndex := 1
	err := reader.NewError(test.NewPosition(1), "ERR1")
	h.RegisterResults(parserIndex, 2, parser.NoCurtailingParsers(), nil, err, parser.EmptyLeftRecCtx())
	cp, rs, actualErr, ok := h.GetResults(parserIndex, 2, parser.EmptyLeftRecCtx())
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Nil(t, rs)
	assert.Equal(t, err, actualErr)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnFalseWhenNoResultWasRegistered(t *testing.T) {
	h := parser.NewHistory()
	parserIndex := 1
	cp, rs, err, ok := h.GetResults(parserIndex, 2, parser.EmptyLeftRecCtx())
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Nil(t, rs)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestRegisterResultShouldHandleMultipleParsers(t *testing.T) {
	h := parser.NewHistory()
	p1Index := 1
	p2Index := 2
	node := ast.NewTerminalNode("t", test.NewPosition(0), nil)
	cp1 := parser.NoCurtailingParsers()
	cp2 := data.NewIntSet(1)
	rs1 := parser.NewResult(node, nil).AsSet()
	var rs2 parser.ResultSet
	err1 := reader.NewError(test.NewPosition(1), "ERR1")
	var err2 reader.Error
	h.RegisterResults(p1Index, 1, cp1, rs1, err1, parser.EmptyLeftRecCtx())
	h.RegisterResults(p2Index, 2, cp2, rs2, err2, parser.EmptyLeftRecCtx())

	actualCP, actualRS, actualErr, ok := h.GetResults(p1Index, 1, parser.EmptyLeftRecCtx())
	assert.Equal(t, cp1, actualCP)
	assert.Equal(t, rs1, actualRS)
	assert.Equal(t, err1, actualErr)
	assert.True(t, ok)

	actualCP, actualRS, actualErr, ok = h.GetResults(p2Index, 2, parser.EmptyLeftRecCtx())
	assert.Equal(t, cp2, actualCP)
	assert.Equal(t, rs2, actualRS)
	assert.Equal(t, err2, actualErr)
	assert.True(t, ok)
}

func TestGetResultsShouldNotReturnCurtailedResult(t *testing.T) {
	h := parser.NewHistory()
	p1Index := 1
	p2Index := 2
	ctx := data.NewIntMap(map[int]int{
		p1Index: 2,
		p2Index: 1,
	})
	cp := data.NewIntSet(p1Index)
	h.RegisterResults(p1Index, 1, cp, nil, nil, ctx)

	ctx = data.NewIntMap(map[int]int{
		p1Index: 1,
		p2Index: 1,
	})
	cp, rs, err, found := h.GetResults(p1Index, 1, ctx)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Nil(t, rs)
	assert.Nil(t, err)
	assert.False(t, found)
}

func TestGetResultsShouldReturnCurtailedResult(t *testing.T) {
	h := parser.NewHistory()
	p1Index := 1
	p2Index := 2
	ctx := data.NewIntMap(map[int]int{
		p1Index: 2,
		p2Index: 1,
	})
	cp := data.NewIntSet(p1Index)
	rs := parser.NewResult(nil, nil).AsSet()
	h.RegisterResults(p1Index, 1, cp, rs, nil, ctx)

	ctx = data.NewIntMap(map[int]int{
		p1Index: 1,
		p2Index: 1,
	})

	ctx = ctx.Inc(p1Index)
	actualCP, actualRS, err, found := h.GetResults(p1Index, 1, ctx)
	assert.Equal(t, cp, actualCP)
	assert.Equal(t, rs, actualRS)
	assert.Nil(t, err)
	assert.True(t, found)
}

func TestResetShouldDeleteData(t *testing.T) {
	h := parser.NewHistory()
	ctx := data.NewIntMap(map[int]int{1: 1})
	cp := data.NewIntSet(1)
	rs := parser.NewResult(nil, nil).AsSet()
	h.RegisterResults(1, 1, cp, rs, nil, ctx)
	h.Reset()
	_, _, _, found := h.GetResults(1, 1, data.EmptyIntMap())
	assert.False(t, found)
}

func TestMemoizeShouldIncreaseLeftRecCtx(t *testing.T) {
	r := test.NewReader(0, 2, false, false)
	h := parser.NewHistory()
	leftRecCtx := parser.EmptyLeftRecCtx()
	parserIndex := 1
	assert.Equal(t, leftRecCtx.Get(parserIndex), 0)

	p := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		assert.Equal(t, leftRecCtx.Get(parserIndex), 1)
		return parser.NoCurtailingParsers(), nil, nil
	})
	h.Memoize(p).Parse(leftRecCtx, r)
}

func TestMemoizeShouldReturnParserResult(t *testing.T) {
	r := test.NewReader(0, 2, false, false)
	h := parser.NewHistory()

	node := ast.NewTerminalNode("a", test.NewPosition(1), "a")
	expectedCP := data.NewIntSet(1)
	expectedRS := parser.NewResult(node, r).AsSet()
	expectedErr := reader.NewError(test.NewPosition(1), "ERR1")

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return expectedCP, expectedRS, expectedErr
	})
	cp, rs, err := h.Memoize(p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, expectedCP, cp)
	assert.Equal(t, expectedRS, rs)
	assert.Equal(t, expectedErr, err)
}

func TestMemoizeShouldRememberResult(t *testing.T) {
	r := test.NewReader(0, 2, false, false)
	h := parser.NewHistory()

	node := ast.NewTerminalNode("a", test.NewPosition(1), "a")
	expectedCP := data.NewIntSet(1)
	expectedRS := parser.NewResult(node, r).AsSet()
	expectedErr := reader.NewError(test.NewPosition(1), "ERR1")

	called := false
	p := h.Memoize(parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		called = true
		return expectedCP, expectedRS, expectedErr
	}))

	p.Parse(parser.EmptyLeftRecCtx(), r)

	called = false

	cp, rs, err := p.Parse(parser.EmptyLeftRecCtx(), r)

	assert.False(t, called, "The parser function should not have been called")
	assert.Equal(t, expectedCP, cp)
	assert.Equal(t, expectedRS, rs)
	assert.Equal(t, expectedErr, err)

	p.Parse(parser.EmptyLeftRecCtx(), test.NewReader(1, 1, false, true))
	assert.True(t, called, "The parser function should have been called for the new position")
}

// The parser calls should be curtailed if the parser was already called more times than the remaining token count plus one
func TestMemoizeShouldCurtailResult(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	h := parser.NewHistory()
	ctx := data.NewIntMap(map[int]int{
		1: 3,
		2: 2,
	})

	called := false
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		called = true
		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "ERR1")
	})
	expectedCP := data.NewIntSet(1)
	cp, rs, err := h.Memoize(p).Parse(ctx, r)
	assert.False(t, called, "The call tree should have been curtailed")
	assert.Equal(t, expectedCP, cp)
	assert.Nil(t, rs)
	assert.Nil(t, err)
}
