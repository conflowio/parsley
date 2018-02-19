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

// Let's define a parser which accepts integer or float numbers.
// The parser would return with all matches, so both 1 and 1.23.
func ExampleAny() {
	number := combinator.Any("number",
		terminal.Integer(),
		terminal.Float(),
	)
	s := parsley.NewSentence(number)
	value, _, _ := s.Evaluate(text.NewReader([]byte("1.23"), "", true), nil)
	fmt.Printf("%T %v\n", value, value)
	// Output: float64 1.23
}

func TestAnyShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := test.NewReader(0, 2, false, false)
	assert.Panics(t, func() { combinator.Any("test").Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r) })
}

func TestAnyShouldHandleOnlyOneParser(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	expectedCP := data.NewIntSet(1)
	expectedRS := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x'), r).AsSet()

	p1 := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return expectedCP, expectedRS, nil
	})

	h := parser.NewHistory()
	cp, rs, err := combinator.Any("test", p1).Parse(h, parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, expectedCP, cp)
	assert.Equal(t, expectedRS, rs)
	assert.Nil(t, err)
	assert.Equal(t, 1, h.CallCount())
}

func TestAnyShouldMergeResults(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	var r1, r2 parser.Result

	p1 := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x'), test.NewReader(0, 2, false, true))
		return parser.NoCurtailingParsers(), r1.AsSet(), nil
	})

	p2 := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		r2 = parser.NewResult(ast.NewTerminalNode("STRING", test.NewPosition(1), 'y'), test.NewReader(1, 1, false, true))
		return data.NewIntSet(1), r2.AsSet(), nil
	})

	p3 := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return data.NewIntSet(2), nil, reader.NewError(test.NewPosition(1), "ERR1")
	})

	p4 := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(2), "ERR2")
	})

	h := parser.NewHistory()
	cp, rs, err := combinator.Any("test", p1, p2, p3, p4).Parse(h, parser.EmptyLeftRecCtx(), r)
	expectedCP := data.NewIntSet(1, 2)
	expectedRS := parser.NewResultSet(r1, r2)
	assert.EqualValues(t, expectedCP, cp)
	assert.EqualValues(t, expectedRS, rs)
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(2), err.Pos())
	assert.Equal(t, "ERR2 at Pos{2}", err.Error())

	assert.Equal(t, 4, h.CallCount())
}

func TestAnyMayReturnEmptyResult(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	p1 := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(2), "ERR1")
	})

	p2 := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(1), "ERR2")
	})

	cp, rs, err := combinator.Any("test", p1, p2).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(2), err.Pos())
	assert.Equal(t, "ERR1 at Pos{2}", err.Error())
}

func TestAnyShouldReturnCustomErrorIfNoParserAdvanced(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	p1 := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(0), "ERR1")
	})

	p2 := parser.Func(func(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return parser.NoCurtailingParsers(), nil, reader.NewError(test.NewPosition(0), "ERR2")
	})

	cp, rs, err := combinator.Any("test", p1, p2).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(0), err.Pos())
	assert.Equal(t, "was expecting test at Pos{0}", err.Error())
}
