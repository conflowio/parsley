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

func TestMemoizeShouldIncreaseleftRecCtx(t *testing.T) {
	r := reader.New([]byte("aa"), true)
	r.ReadRune()
	h := parser.NewHistory()
	leftRecCtx := parser.EmptyLeftRecCtx()
	parserIndex := h.GetParserIndex("p1")
	assert.Equal(t, leftRecCtx.Get(parserIndex), 0)

	p := parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		assert.Equal(t, leftRecCtx.Get(parserIndex), 1)
		return nil
	})
	combinator.Memoize("p1", h, p).Parse(leftRecCtx, r)
}

func TestMemoizeShouldReturnParserResult(t *testing.T) {
	r := reader.New([]byte("a"), true)
	h := parser.NewHistory()
	ctx := parser.EmptyLeftRecCtx()

	node := ast.NewTerminalNode("a", reader.NewPosition(0, 1, 2), "a")
	expected := parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(node, r))

	p := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return expected
	})
	actual := combinator.Memoize("p1", h, p).Parse(ctx, r)
	assert.Equal(t, expected, actual)
}

func TestMemoizeShouldRememberResult(t *testing.T) {
	r := reader.New([]byte("a"), true)
	h := parser.NewHistory()
	ctx := parser.EmptyLeftRecCtx()

	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	expected := parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(node, r))

	called := false
	p := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		called = true
		return expected
	})
	combinator.Memoize("p1", h, p).Parse(ctx, r)

	called = false

	actual := combinator.Memoize("p1", h, p).Parse(ctx, r)

	assert.False(t, called, "The parser function should not have been called")
	assert.Equal(t, expected, actual)

	r.ReadRune()

	combinator.Memoize("p1", h, p).Parse(ctx, r)
	assert.True(t, called, "The parser function should have been called for the new position")
}

// The parser calls should be curtailed if the parser was already called more times than the remaining token count plus one
func TestMemoizeShouldCurtailResult(t *testing.T) {
	r := reader.New([]byte("a"), true)
	h := parser.NewHistory()
	ctx := parser.EmptyLeftRecCtx().
		Inc(h.GetParserIndex("p1")).
		Inc(h.GetParserIndex("p1")).
		Inc(h.GetParserIndex("p1")).
		Inc(h.GetParserIndex("p2"))

	called := false
	p := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		called = true
		return nil
	})
	expected := parser.NewParserResult(data.NewIntSet(h.GetParserIndex("p1")))
	results := combinator.Memoize("p1", h, p).Parse(ctx, r)
	assert.False(t, called, "The call tree should have been curtailed")
	assert.Equal(t, expected, results)
}
