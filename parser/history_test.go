package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
)

func TestRegisterResultShouldSaveResultForPosition(t *testing.T) {
	h := parser.NewHistory()
	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	results := parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(node, nil))
	h.RegisterResults(h.GetParserIndex("p1"), 2, results, parser.EmptyLeftRecCtx())

	actual, ok := h.GetResults(h.GetParserIndex("p1"), 2, parser.EmptyLeftRecCtx())
	assert.Equal(t, results, actual)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnNilResult(t *testing.T) {
	h := parser.NewHistory()
	h.RegisterResults(h.GetParserIndex("p1"), 2, nil, parser.EmptyLeftRecCtx())
	results, ok := h.GetResults(h.GetParserIndex("p1"), 2, parser.EmptyLeftRecCtx())
	assert.Nil(t, results)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnFalseWhenNoResultWasRegistered(t *testing.T) {
	h := parser.NewHistory()
	results, ok := h.GetResults(h.GetParserIndex("p1"), 2, parser.EmptyLeftRecCtx())
	assert.Nil(t, results)
	assert.False(t, ok)
}

func TestRegisterResultShouldHandleMultipleParsers(t *testing.T) {
	h := parser.NewHistory()
	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	results := parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(node, nil))
	h.RegisterResults(h.GetParserIndex("p1"), 1, results, parser.EmptyLeftRecCtx())
	h.RegisterResults(h.GetParserIndex("p2"), 2, nil, parser.EmptyLeftRecCtx())

	actual, ok := h.GetResults(h.GetParserIndex("p1"), 1, parser.EmptyLeftRecCtx())
	assert.Equal(t, results, actual)
	assert.True(t, ok)

	results, ok = h.GetResults(h.GetParserIndex("p2"), 2, parser.EmptyLeftRecCtx())
	assert.Nil(t, results)
	assert.True(t, ok)
}

func TestGetResultsShouldNotReturnCurtailedResult(t *testing.T) {
	h := parser.NewHistory()
	ctx := data.NewIntMap(map[int]int{
		h.GetParserIndex("p1"): 2,
		h.GetParserIndex("p2"): 1,
	})
	curtailingParsers := data.NewIntSet(h.GetParserIndex("p1"))
	h.RegisterResults(h.GetParserIndex("p1"), 1, parser.NewParserResult(curtailingParsers), ctx)

	ctx = data.NewIntMap(map[int]int{
		h.GetParserIndex("p1"): 1,
		h.GetParserIndex("p2"): 1,
	})
	results, found := h.GetResults(h.GetParserIndex("p1"), 1, ctx)
	assert.Nil(t, results)
	assert.False(t, found)
}

func TestGetResultsShouldReturnCurtailedResult(t *testing.T) {
	h := parser.NewHistory()
	ctx := data.NewIntMap(map[int]int{
		h.GetParserIndex("p1"): 2,
		h.GetParserIndex("p2"): 1,
	})
	curtailingParsers := data.NewIntSet(h.GetParserIndex("p1"))
	cResults := parser.NewParserResult(curtailingParsers, parser.NewResult(nil, nil))
	h.RegisterResults(h.GetParserIndex("p1"), 1, cResults, ctx)

	ctx = data.NewIntMap(map[int]int{
		h.GetParserIndex("p1"): 1,
		h.GetParserIndex("p2"): 1,
	})

	ctx = ctx.Inc(h.GetParserIndex("p1"))
	results, found := h.GetResults(h.GetParserIndex("p1"), 1, ctx)
	assert.Equal(t, cResults, results)
	assert.True(t, found)
}
