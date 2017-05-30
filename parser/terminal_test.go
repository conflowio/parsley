package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
)

func assertCursor(t *testing.T, pos int, line int, col int, r *reader.Reader) {
	assert.Equal(t, reader.NewPosition(pos, line, col), r.Cursor())
}

func TestEmptyWillAlwaysReturnWithResult(t *testing.T) {
	r := reader.New([]byte{}, true)
	res := parser.Empty()(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode(ast.EMPTY, reader.NewPosition(0, 1, 1), nil)
	assert.Equal(t, parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(expectedNode, r)), res)
	assertCursor(t, 0, 1, 1, r)
}

func TestEndShouldMatchEmptyReader(t *testing.T) {
	r := reader.New([]byte{}, true)
	res := parser.End()(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode(ast.EOF, reader.NewPosition(0, 1, 1), nil)
	assert.Equal(t, parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(expectedNode, r)), res)
	assertCursor(t, 0, 1, 1, r)
}

func TestEndShouldMatchEOF(t *testing.T) {
	r := reader.New([]byte("a"), true)
	r.ReadRune()
	res := parser.End()(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode(ast.EOF, reader.NewPosition(1, 1, 2), nil)
	assert.Equal(t, parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(expectedNode, r)), res)
	assertCursor(t, 1, 1, 2, r)
}

func TestEndShouldNotMatchNotEOF(t *testing.T) {
	r := reader.New([]byte("a"), true)
	res := parser.End()(parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	assertCursor(t, 0, 1, 1, r)
}

func TestRuneShouldMatchCharacter(t *testing.T) {
	r := reader.New([]byte("a"), true)
	res := parser.Rune('a', "A")(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("A", reader.NewPosition(0, 1, 1), 'a')
	assert.Equal(t, parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(expectedNode, r)), res)
	assertCursor(t, 1, 1, 2, r)
}

func TestRuneShouldNotUseSpecialChars(t *testing.T) {
	r := reader.New([]byte("a"), true)
	res := parser.Rune('.', ".")(parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	assertCursor(t, 0, 1, 1, r)
}

func TestRuneShouldMatchUnicodeCharacter(t *testing.T) {
	r := reader.New([]byte("🍕"), true)
	res := parser.Rune('🍕', "PIZZA")(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("PIZZA", reader.NewPosition(0, 1, 1), '🍕')
	assert.Equal(t, parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(expectedNode, r)), res)
	assertCursor(t, 4, 1, 2, r)
}

func TestRuneShouldNotMatchCharacter(t *testing.T) {
	r := reader.New([]byte("a"), true)
	res := parser.Rune('b', "B")(parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	assertCursor(t, 0, 1, 1, r)
}
