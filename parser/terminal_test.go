package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
)

func assertCursor(t *testing.T, pos int, line int, col int, r *reader.Reader) {
	assert.Equal(t, reader.NewPosition(pos, line, col), r.Cursor())
}

func TestEndShouldMatchEmptyReader(t *testing.T) {
	r := reader.New([]byte{}, true)
	res := parser.End()(data.NewIntMap(nil), r)
	expectedNode := ast.NewTerminalNode(reader.EOF, reader.NewPosition(0, 1, 1), nil)
	assert.Equal(t, parser.NewParserResult(data.NewIntSet(), parser.NewResult(expectedNode, r)), res)
	assertCursor(t, 0, 1, 1, r)
}

func TestEndShouldMatchEOF(t *testing.T) {
	r := reader.New([]byte("a"), true)
	r.ReadRune()
	res := parser.End()(data.NewIntMap(nil), r)
	expectedNode := ast.NewTerminalNode(reader.EOF, reader.NewPosition(1, 1, 2), nil)
	assert.Equal(t, parser.NewParserResult(data.NewIntSet(), parser.NewResult(expectedNode, r)), res)
	assertCursor(t, 1, 1, 2, r)
}

func TestEndShouldNotMatchNotEOF(t *testing.T) {
	r := reader.New([]byte("a"), true)
	res := parser.End()(data.NewIntMap(nil), r)
	assert.Nil(t, res)
	assertCursor(t, 0, 1, 1, r)
}

func TestRuneShouldMatchCharacter(t *testing.T) {
	r := reader.New([]byte("a"), true)
	res := parser.Rune('a', "A")(data.NewIntMap(nil), r)
	expectedNode := ast.NewTerminalNode("A", reader.NewPosition(0, 1, 1), 'a')
	assert.Equal(t, parser.NewParserResult(data.NewIntSet(), parser.NewResult(expectedNode, r)), res)
	assertCursor(t, 1, 1, 2, r)
}

func TestRuneShouldNotUseSpecialChars(t *testing.T) {
	r := reader.New([]byte("a"), true)
	res := parser.Rune('.', ".")(data.NewIntMap(nil), r)
	assert.Nil(t, res)
	assertCursor(t, 0, 1, 1, r)
}

func TestRuneShouldMatchUnicodeCharacter(t *testing.T) {
	r := reader.New([]byte("🍕"), true)
	res := parser.Rune('🍕', "PIZZA")(data.NewIntMap(nil), r)
	expectedNode := ast.NewTerminalNode("PIZZA", reader.NewPosition(0, 1, 1), '🍕')
	assert.Equal(t, parser.NewParserResult(data.NewIntSet(), parser.NewResult(expectedNode, r)), res)
	assertCursor(t, 4, 1, 2, r)
}

func TestRuneShouldNotMatchCharacter(t *testing.T) {
	r := reader.New([]byte("a"), true)
	res := parser.Rune('b', "B")(data.NewIntMap(nil), r)
	assert.Nil(t, res)
	assertCursor(t, 0, 1, 1, r)
}
