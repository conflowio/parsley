package combinator_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
)

func TestOrShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	ctx := parser.EmptyLeftRecCtx()
	assert.Panics(t, func() { combinator.Or().Parse(ctx, r) })
}

func TestOrShouldHandleOnlyOneParser(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	ctx := parser.EmptyLeftRecCtx()

	var r1 parser.Result

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r)
		return parser.NewParserResult(parser.NoCurtailingParsers(), r1)
	})

	results := combinator.Or(p1).Parse(ctx, r)
	assert.Equal(t, parser.NewParserResult(parser.NoCurtailingParsers(), r1), results)
}

func TestOrShouldMergeResults(t *testing.T) {
	parser.Stat.Reset()
	r := reader.New([]byte("ab"), true)
	ctx := parser.EmptyLeftRecCtx()

	var r1, r2 parser.Result

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r)
		return parser.NewParserResult(parser.NoCurtailingParsers(), r1)
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch1, _, _ := r.ReadRune()
		ch2, _, _ := r.ReadRune()
		r2 = parser.NewResult(ast.NewTerminalNode("STRING", pos, string([]rune{ch1, ch2})), r)
		return parser.NewParserResult(data.NewIntSet(1), r2)
	})

	p3 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return parser.NewParserResult(data.NewIntSet(2))
	})

	p4 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return nil
	})

	results := combinator.Or(p1, p2, p3, p4).Parse(ctx, r)
	curtailingParsers := data.NewIntSet(1, 2)
	assert.EqualValues(t, parser.NewParserResult(curtailingParsers, r1, r2), results)

	assert.Equal(t, 4, parser.Stat.GetSumCallCount())
}

func TestOrMayReturnEmptyResult(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	ctx := parser.EmptyLeftRecCtx()

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return nil
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return nil
	})

	results := combinator.Or(p1, p2).Parse(ctx, r)
	assert.Equal(t, parser.NewParserResult(parser.NoCurtailingParsers()), results)
}

func TestOrShouldCloneReadersForAllParsers(t *testing.T) {
	parser.Stat.Reset()
	r := reader.New([]byte("ab"), true)
	ctx := parser.EmptyLeftRecCtx()

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		assert.Equal(t, reader.NewPosition(0, 1, 1), r.Cursor())
		r.ReadRune()
		return nil
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		assert.Equal(t, reader.NewPosition(0, 1, 1), r.Cursor())
		r.ReadRune()
		return nil
	})

	combinator.Or(p1, p2).Parse(ctx, r)
	assert.Equal(t, 2, parser.Stat.GetSumCallCount())
}
