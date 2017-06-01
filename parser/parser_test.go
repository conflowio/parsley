package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
)

type testPosition struct {
	pos int
}

func (t testPosition) Pos() int       { return t.pos }
func (t testPosition) String() string { return "testPosition" }

type testReader struct {
	cur       testPosition
	remaining int
	isEOF     bool
}

func (t testReader) Clone() reader.Reader    { return t }
func (t testReader) Cursor() reader.Position { return t.cur }
func (t testReader) Remaining() int          { return t.remaining }
func (t testReader) IsEOF() bool             { return t.isEOF }
func (t testReader) String() string          { return "testReader" }

func newTestReader(pos int, remaining int, isEOF bool) testReader {
	return testReader{testPosition{pos}, remaining, isEOF}
}

func TesParserFuncShouldCallFunction(t *testing.T) {
	expectedLeftRecCtx := data.NewIntMap(map[int]int{1: 2})
	expectedReader := newTestReader(0, 1, false)
	expectedCurtailingParsers := data.NewIntSet(1)
	expectedResultSet := parser.NewResult(nil, newTestReader(1, 0, false)).AsSet()
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
