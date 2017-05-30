package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
)

func TestRegisterCallShouldIncreaseSumCallCount(t *testing.T) {
	h := parser.NewHistory()
	h.RegisterCall()
	assert.Equal(t, h.GetSumCallCount(), 2)
	h.RegisterCall()
	assert.Equal(t, h.GetSumCallCount(), 3)
}

func TestRegisterResultShouldSaveResultForPosition(t *testing.T) {
	h := parser.NewHistory()
	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	results := parser.NewParserResult(data.NewIntSet(), parser.NewResult(node, nil))
	h.RegisterResults(h.GetParserIndex("p1"), 2, results, data.NewIntMap(nil))

	actual, ok := h.GetResults(h.GetParserIndex("p1"), 2, data.NewIntMap(nil))
	assert.Equal(t, results, actual)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnNilResult(t *testing.T) {
	h := parser.NewHistory()
	h.RegisterResults(h.GetParserIndex("p1"), 2, nil, data.NewIntMap(nil))
	results, ok := h.GetResults(h.GetParserIndex("p1"), 2, data.NewIntMap(nil))
	assert.Nil(t, results)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnFalseWhenNoResultWasRegistered(t *testing.T) {
	h := parser.NewHistory()
	results, ok := h.GetResults(h.GetParserIndex("p1"), 2, data.NewIntMap(nil))
	assert.Nil(t, results)
	assert.False(t, ok)
}

func TestRegisterResultShouldHandleMultipleParsers(t *testing.T) {
	h := parser.NewHistory()
	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	results := parser.NewParserResult(data.NewIntSet(), parser.NewResult(node, nil))
	h.RegisterResults(h.GetParserIndex("p1"), 1, results, data.NewIntMap(nil))
	h.RegisterResults(h.GetParserIndex("p2"), 2, nil, data.NewIntMap(nil))

	actual, ok := h.GetResults(h.GetParserIndex("p1"), 1, data.NewIntMap(nil))
	assert.Equal(t, results, actual)
	assert.True(t, ok)

	results, ok = h.GetResults(h.GetParserIndex("p2"), 2, data.NewIntMap(nil))
	assert.Nil(t, results)
	assert.True(t, ok)
}

func TestGetSumCallCountShouldStartAtOne(t *testing.T) {
	h := parser.NewHistory()
	assert.Equal(t, h.GetSumCallCount(), 1)
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
