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

func TestMemoizeShouldIncreaseLeftRecCtx(t *testing.T) {
	r := test.NewReader(0, 2, false, false)
	h := parser.NewHistory()
	leftRecCtx := parser.EmptyLeftRecCtx()
	parserIndex := h.GetParserIndex("p1")
	assert.Equal(t, leftRecCtx.Get(parserIndex), 0)

	p := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		assert.Equal(t, leftRecCtx.Get(parserIndex), 1)
		return parser.NoCurtailingParsers(), nil, nil
	})
	combinator.Memoize("p1", h, p).Parse(leftRecCtx, r)
}

func TestMemoizeShouldReturnParserResult(t *testing.T) {
	r := test.NewReader(0, 2, false, false)
	h := parser.NewHistory()

	node := ast.NewTerminalNode("a", test.NewPosition(1), "a")
	expectedCP := data.NewIntSet(1)
	expectedRS := parser.NewResult(node, r).AsSet()
	expectedErr := parser.NewError(test.NewPosition(1), "ERR1")

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return expectedCP, expectedRS, expectedErr
	})
	cp, rs, err := combinator.Memoize("p1", h, p).Parse(parser.EmptyLeftRecCtx(), r)
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
	expectedErr := parser.NewError(test.NewPosition(1), "ERR1")

	called := false
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		called = true
		return expectedCP, expectedRS, expectedErr
	})
	combinator.Memoize("p1", h, p).Parse(parser.EmptyLeftRecCtx(), r)

	called = false

	cp, rs, err := combinator.Memoize("p1", h, p).Parse(parser.EmptyLeftRecCtx(), r)

	assert.False(t, called, "The parser function should not have been called")
	assert.Equal(t, expectedCP, cp)
	assert.Equal(t, expectedRS, rs)
	assert.Equal(t, expectedErr, err)

	combinator.Memoize("p1", h, p).Parse(parser.EmptyLeftRecCtx(), test.NewReader(1, 1, false, true))
	assert.True(t, called, "The parser function should have been called for the new position")
}

// The parser calls should be curtailed if the parser was already called more times than the remaining token count plus one
func TestMemoizeShouldCurtailResult(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	h := parser.NewHistory()
	ctx := data.NewIntMap(map[int]int{
		h.GetParserIndex("p1"): 3,
		h.GetParserIndex("p2"): 2,
	})

	called := false
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		called = true
		return parser.NoCurtailingParsers(), nil, parser.NewError(test.NewPosition(1), "ERR1")
	})
	expectedCP := data.NewIntSet(h.GetParserIndex("p1"))
	cp, rs, err := combinator.Memoize("p1", h, p).Parse(ctx, r)
	assert.False(t, called, "The call tree should have been curtailed")
	assert.Equal(t, expectedCP, cp)
	assert.Nil(t, rs)
	assert.Nil(t, err)
}
