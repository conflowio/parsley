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

func TestAnyShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := test.NewReader(0, 2, false, false)
	assert.Panics(t, func() { combinator.Any().Parse(parser.EmptyLeftRecCtx(), r) })
}

func TestAnyShouldHandleOnlyOneParser(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	expectedCP := data.NewIntSet(1)
	expectedRS := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x'), r).AsSet()

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return expectedCP, expectedRS
	})

	cp, rs := combinator.Any(p1).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, expectedCP, cp)
	assert.Equal(t, expectedRS, rs)
}

func TestAnyShouldMergeResults(t *testing.T) {
	parser.Stat.Reset()
	r := test.NewReader(0, 2, false, false)

	var r1, r2 parser.Result

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x'), test.NewReader(0, 2, false, true))
		return parser.NoCurtailingParsers(), r1.AsSet()
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		r2 = parser.NewResult(ast.NewTerminalNode("STRING", test.NewPosition(1), 'y'), test.NewReader(1, 1, false, true))
		return data.NewIntSet(1), r2.AsSet()
	})

	p3 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return data.NewIntSet(2), nil
	})

	p4 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})

	cp, rs := combinator.Any(p1, p2, p3, p4).Parse(parser.EmptyLeftRecCtx(), r)
	expectedCP := data.NewIntSet(1, 2)
	expectedRS := parser.NewResultSet(r1, r2)
	assert.EqualValues(t, expectedCP, cp)
	assert.EqualValues(t, expectedRS, rs)

	assert.Equal(t, 4, parser.Stat.GetSumCallCount())
}

func TestAnyMayReturnEmptyResult(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})

	cp, rs := combinator.Any(p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
}
