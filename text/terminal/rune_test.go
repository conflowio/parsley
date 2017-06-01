package terminal_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
)

func assertCursor(t *testing.T, pos int, line int, col int, r *text.Reader) {
	assert.Equal(t, text.NewPosition(pos, line, col), r.Cursor())
}

func TestRuneShouldMatchCharacter(t *testing.T) {
	r := text.NewReader([]byte("a"), true)
	_, res := terminal.Rune('a', "A")(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("A", text.NewPosition(0, 1, 1), 'a')
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assertCursor(t, 1, 1, 2, r)
}

func TestRuneShouldNotUseSpecialChars(t *testing.T) {
	r := text.NewReader([]byte("a"), true)
	_, res := terminal.Rune('.', ".")(parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	assertCursor(t, 0, 1, 1, r)
}

func TestRuneShouldMatchUnicodeCharacter(t *testing.T) {
	r := text.NewReader([]byte("🍕"), true)
	_, res := terminal.Rune('🍕', "PIZZA")(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("PIZZA", text.NewPosition(0, 1, 1), '🍕')
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assertCursor(t, 4, 1, 2, r)
}

func TestRuneShouldNotMatchCharacter(t *testing.T) {
	r := text.NewReader([]byte("a"), true)
	_, res := terminal.Rune('b', "B")(parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	assertCursor(t, 0, 1, 1, r)
}
