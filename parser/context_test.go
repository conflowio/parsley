package parser_test

import (
	"testing"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/parser"
	"github.com/opsidian/parsec/reader"
	"github.com/stretchr/testify/assert"
)

func TestPushShouldIncreaseCallCount(t *testing.T) {
	c := parser.NewContext()
	c.Push("p1", 2)
	assert.Equal(t, c.GetCalls("p1", 2), 1)
	c.Push("p1", 2)
	assert.Equal(t, c.GetCalls("p1", 2), 2)
}

func TestPopShouldDecreaseCallCount(t *testing.T) {
	c := parser.NewContext()
	c.Push("p1", 2)
	c.Push("p1", 2)
	assert.Equal(t, c.GetCalls("p1", 2), 2)
	c.Pop("p1", 2)
	assert.Equal(t, c.GetCalls("p1", 2), 1)
	c.Pop("p1", 2)
	assert.Equal(t, c.GetCalls("p1", 2), 0)
}

func TestPopShouldPanicIfNoPushWasCalled(t *testing.T) {
	c := parser.NewContext()
	assert.Panics(t, func() { c.Pop("p1", 2) })
}

func TestPushPopShouldHandleMultipleParsers(t *testing.T) {
	c := parser.NewContext()
	c.Push("p1", 1)
	c.Push("p2", 2)
	assert.Equal(t, c.GetCalls("p1", 1), 1)
	assert.Equal(t, c.GetCalls("p2", 2), 1)
	c.Pop("p1", 1)
	assert.Equal(t, c.GetCalls("p1", 1), 0)
	c.Pop("p2", 2)
	assert.Equal(t, c.GetCalls("p2", 2), 0)
}

func TestRegiterCallShouldIncreaseSumCallCount(t *testing.T) {
	c := parser.NewContext()
	c.RegisterCall()
	assert.Equal(t, c.GetSumCallCount(), 2)
	c.RegisterCall()
	assert.Equal(t, c.GetSumCallCount(), 3)
}

func TestRegisterResultShouldSaveResultForPosition(t *testing.T) {
	c := parser.NewContext()
	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	results := parser.NewResult(node, nil).AsList()
	c.RegisterResults("p1", 2, results)

	actual, ok := c.GetResults("p1", 2)
	assert.Equal(t, results, actual)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnNilResult(t *testing.T) {
	c := parser.NewContext()
	c.RegisterResults("p1", 2, nil)
	results, ok := c.GetResults("p1", 2)
	assert.Nil(t, results)
	assert.True(t, ok)
}

func TestRegisterResultShouldReturnFalseWhenNoResultWasRegistered(t *testing.T) {
	c := parser.NewContext()
	results, ok := c.GetResults("p1", 2)
	assert.Nil(t, results)
	assert.False(t, ok)
}

func TestRegisterResultShouldHandleMultipleParsers(t *testing.T) {
	c := parser.NewContext()
	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	results := parser.NewResult(node, nil).AsList()
	c.RegisterResults("p1", 1, results)
	c.RegisterResults("p2", 2, nil)

	actual, ok := c.GetResults("p1", 1)
	assert.Equal(t, results, actual)
	assert.True(t, ok)

	results, ok = c.GetResults("p2", 2)
	assert.Nil(t, results)
	assert.True(t, ok)
}

func TestGetSumCallCountShouldStartAtOne(t *testing.T) {
	c := parser.NewContext()
	assert.Equal(t, c.GetSumCallCount(), 1)
}

func TestGetResultsShouldNotReturnCurtailedResult(t *testing.T) {
	c := parser.NewContext()
	c.Push("p1", 1)
	c.Push("p2", 1)
	c.Push("p1", 1)
	c.RegisterResults("p1", 1, c.NewCurtailedResults(1))

	c.Pop("p1", 1)
	c.Pop("p2", 1)
	c.Pop("p1", 1)

	c.Push("p3", 1)
	c.Push("p1", 1)
	c.Push("p2", 1)
	results, found := c.GetResults("p1", 1)
	assert.Nil(t, results)
	assert.False(t, found)
}

func TestGetResultsShouldReturnCurtailedResult(t *testing.T) {
	c := parser.NewContext()
	c.Push("p1", 1)
	c.Push("p2", 1)
	c.Push("p1", 1)
	cResults := c.NewCurtailedResults(1)
	c.Pop("p1", 1)
	c.Pop("p2", 1)
	c.Pop("p1", 1)
	c.RegisterResults("p1", 1, cResults)
	c.RegisterResults("p2", 1, cResults)

	var found bool
	c.Push("p3", 1)
	_, found = c.GetResults("p3", 1)
	assert.False(t, found)

	c.Push("p1", 1)
	_, found = c.GetResults("p1", 1)
	assert.False(t, found)

	c.Push("p2", 1)
	_, found = c.GetResults("p2", 1)
	assert.False(t, found)

	c.Push("p1", 1)
	results, found := c.GetResults("p1", 1)
	assert.Equal(t, cResults, results)
	assert.True(t, found)
}
