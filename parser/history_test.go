// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestRegisterResultShouldSaveResultForPosition(t *testing.T) {
	h := parser.NewHistory()
	node := ast.NewTerminalNode("t", test.NewPosition(0), nil)
	cp := parser.NoCurtailingParsers()
	rs := parser.NewResult(node, nil).AsSet()
	err := parser.NewError(test.NewPosition(1), "ERR1")
	h.RegisterResults(h.GetParserIndex("p1"), 2, cp, rs, err, parser.EmptyLeftRecCtx())

	actualCP, actualRS, actualErr, ok := h.GetResults(h.GetParserIndex("p1"), 2, parser.EmptyLeftRecCtx())
	assert.Equal(t, cp, actualCP)
	assert.Equal(t, rs, actualRS)
	assert.Equal(t, err, actualErr)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnNilResult(t *testing.T) {
	h := parser.NewHistory()
	err := parser.NewError(test.NewPosition(1), "ERR1")
	h.RegisterResults(h.GetParserIndex("p1"), 2, parser.NoCurtailingParsers(), nil, err, parser.EmptyLeftRecCtx())
	cp, rs, actualErr, ok := h.GetResults(h.GetParserIndex("p1"), 2, parser.EmptyLeftRecCtx())
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Nil(t, rs)
	assert.Equal(t, err, actualErr)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnFalseWhenNoResultWasRegistered(t *testing.T) {
	h := parser.NewHistory()
	cp, rs, err, ok := h.GetResults(h.GetParserIndex("p1"), 2, parser.EmptyLeftRecCtx())
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Nil(t, rs)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestRegisterResultShouldHandleMultipleParsers(t *testing.T) {
	h := parser.NewHistory()
	node := ast.NewTerminalNode("t", test.NewPosition(0), nil)
	cp1 := parser.NoCurtailingParsers()
	cp2 := data.NewIntSet(1)
	rs1 := parser.NewResult(node, nil).AsSet()
	var rs2 parser.ResultSet
	err1 := parser.NewError(test.NewPosition(1), "ERR1")
	var err2 parser.Error
	h.RegisterResults(h.GetParserIndex("p1"), 1, cp1, rs1, err1, parser.EmptyLeftRecCtx())
	h.RegisterResults(h.GetParserIndex("p2"), 2, cp2, rs2, err2, parser.EmptyLeftRecCtx())

	actualCP, actualRS, actualErr, ok := h.GetResults(h.GetParserIndex("p1"), 1, parser.EmptyLeftRecCtx())
	assert.Equal(t, cp1, actualCP)
	assert.Equal(t, rs1, actualRS)
	assert.Equal(t, err1, actualErr)
	assert.True(t, ok)

	actualCP, actualRS, actualErr, ok = h.GetResults(h.GetParserIndex("p2"), 2, parser.EmptyLeftRecCtx())
	assert.Equal(t, cp2, actualCP)
	assert.Equal(t, rs2, actualRS)
	assert.Equal(t, err2, actualErr)
	assert.True(t, ok)
}

func TestGetResultsShouldNotReturnCurtailedResult(t *testing.T) {
	h := parser.NewHistory()
	ctx := data.NewIntMap(map[int]int{
		h.GetParserIndex("p1"): 2,
		h.GetParserIndex("p2"): 1,
	})
	cp := data.NewIntSet(h.GetParserIndex("p1"))
	h.RegisterResults(h.GetParserIndex("p1"), 1, cp, nil, nil, ctx)

	ctx = data.NewIntMap(map[int]int{
		h.GetParserIndex("p1"): 1,
		h.GetParserIndex("p2"): 1,
	})
	cp, rs, err, found := h.GetResults(h.GetParserIndex("p1"), 1, ctx)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Nil(t, rs)
	assert.Nil(t, err)
	assert.False(t, found)
}

func TestGetResultsShouldReturnCurtailedResult(t *testing.T) {
	h := parser.NewHistory()
	ctx := data.NewIntMap(map[int]int{
		h.GetParserIndex("p1"): 2,
		h.GetParserIndex("p2"): 1,
	})
	cp := data.NewIntSet(h.GetParserIndex("p1"))
	rs := parser.NewResult(nil, nil).AsSet()
	h.RegisterResults(h.GetParserIndex("p1"), 1, cp, rs, nil, ctx)

	ctx = data.NewIntMap(map[int]int{
		h.GetParserIndex("p1"): 1,
		h.GetParserIndex("p2"): 1,
	})

	ctx = ctx.Inc(h.GetParserIndex("p1"))
	actualCP, actualRS, err, found := h.GetResults(h.GetParserIndex("p1"), 1, ctx)
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
