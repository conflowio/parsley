package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestParserFuncShouldCallFunction(t *testing.T) {
	expectedLeftRecCtx := data.NewIntMap(map[int]int{1: 2})
	expectedReader := test.NewReader(0, 1, false, false)
	expectedCurtailingParsers := data.NewIntSet(1)
	expectedResultSet := parser.NewResult(nil, test.NewReader(1, 0, false, false)).AsSet()
	var actualLeftRecCtx data.IntMap
	var actualReader reader.Reader
	parserFunc := parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		actualLeftRecCtx = leftRecCtx
		actualReader = r
		return expectedCurtailingParsers, expectedResultSet
	})

	actualCurtailingParsers, actualResultSet := parserFunc.Parse(expectedLeftRecCtx, expectedReader)

	assert.Equal(t, expectedCurtailingParsers, actualCurtailingParsers)
	assert.Equal(t, expectedResultSet, actualResultSet)
	assert.Equal(t, expectedLeftRecCtx, actualLeftRecCtx)
	assert.Equal(t, expectedReader, actualReader)
}
