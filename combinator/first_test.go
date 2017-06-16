package combinator_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestFirstShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := test.NewReader(0, 2, false, false)
	assert.Panics(t, func() { combinator.First().Parse(parser.EmptyLeftRecCtx(), r) })
}

func TestFirstShouldHandleOnlyOneParser(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	expectedCP := data.NewIntSet(1)
	expectedRS := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x'), r).AsSet()

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return expectedCP, expectedRS
	})

	cp, rs := combinator.First(p1).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, expectedCP, cp)
	assert.Equal(t, expectedRS, rs)
}

func TestFirstShouldMergeCurtailingParsers(t *testing.T) {
	parser.Stat.Reset()
	r := test.NewReader(0, 2, false, false)

	var res parser.Result

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return data.NewIntSet(1), parser.NewResultSet()
	})

	p3 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		res = parser.NewResult(ast.NewTerminalNode("STRING", test.NewPosition(0), "TEST"), test.NewReader(1, 1, false, true))
		return data.NewIntSet(2), res.AsSet()
	})

	cp, rs := combinator.First(p1, p2, p3).Parse(parser.EmptyLeftRecCtx(), r)
	expectedCP := data.NewIntSet(1, 2)
	expectedRS := res.AsSet()
	assert.EqualValues(t, expectedCP, cp)
	assert.EqualValues(t, expectedRS, rs)

	assert.Equal(t, 3, parser.Stat.GetSumCallCount())
}

func TestFirstShouldStopAtFirstResult(t *testing.T) {
	parser.Stat.Reset()
	r := test.NewReader(0, 2, false, false)

	var res1, res2 parser.Result

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		res1 = parser.NewResult(ast.NewTerminalNode("STRING", test.NewPosition(0), "TEST"), test.NewReader(1, 1, false, true))
		return parser.NoCurtailingParsers(), res1.AsSet()
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		res2 = parser.NewResult(ast.NewTerminalNode("STRING", test.NewPosition(0), "TEST2"), test.NewReader(1, 1, false, true))
		return parser.NoCurtailingParsers(), res2.AsSet()
	})

	_, rs := combinator.First(p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, res1.AsSet(), rs)

	assert.Equal(t, 1, parser.Stat.GetSumCallCount())
}

func TestFirstMayReturnEmptyResult(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})

	cp, rs := combinator.First(p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
}
