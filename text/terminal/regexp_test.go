package terminal_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
)

func TestRegexpShouldMatchPattern(t *testing.T) {
	r := text.NewReader([]byte("abc123"), true)
	_, res := terminal.Regexp("[a-z]+", false, 0, "TEST")(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("TEST", text.NewPosition(0, 1, 1), "abc")
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assert.Equal(t, text.NewPosition(3, 1, 4), r.Cursor())
}

func TestRegexpShouldIgnoreWhitespaces(t *testing.T) {
	r := text.NewReader([]byte("x   abc123"), true)
	r.ReadRune()
	_, res := terminal.Regexp("[a-z]+", false, 0, "TEST")(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("TEST", text.NewPosition(4, 1, 5), "abc")
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assert.Equal(t, text.NewPosition(7, 1, 8), r.Cursor())
}

func TestRegexpShouldNotIgnoreWhitespaces(t *testing.T) {
	r := text.NewReader([]byte("x   abc123a"), true)
	r.ReadRune()
	_, res := terminal.Regexp("[ a-z]+", true, 0, "TEST")(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("TEST", text.NewPosition(1, 1, 2), "   abc")
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assert.Equal(t, text.NewPosition(7, 1, 8), r.Cursor())
}

func TestRegexpShouldUseMatchIfTokenEmpty(t *testing.T) {
	r := text.NewReader([]byte("abc123"), true)
	_, res := terminal.Regexp("[a-z]+", false, 0, "")(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("abc", text.NewPosition(0, 1, 1), "abc")
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assert.Equal(t, text.NewPosition(3, 1, 4), r.Cursor())
}

func TestRegexpNotMatchingPattern(t *testing.T) {
	r := text.NewReader([]byte("abc123"), true)
	_, res := terminal.Regexp("[0-9]+", false, 0, "TEST")(parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
}

func TestRegexpShouldPanicIfInvalidGroupIndex(t *testing.T) {
	r := text.NewReader([]byte("abc123"), true)
	assert.Panics(t, func() { terminal.Regexp("[a-z]+", false, 1, "TEST")(parser.EmptyLeftRecCtx(), r) })
}
