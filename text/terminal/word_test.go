package terminal_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
)

func TestWordShouldMatchExactString(t *testing.T) {
	r := text.NewReader([]byte("foo bar"), true)
	_, res := terminal.Word("foo", "TEST", 5).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, res, parser.NewResult(ast.NewTerminalNode("TEST", text.NewPosition(0, 1, 1), 5), r).AsSet())
}

func TestWordShouldNotMatchWordPartially(t *testing.T) {
	r := text.NewReader([]byte("foobar"), true)
	_, res := terminal.Word("foo", "TEST", 5).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
}

func TestWordShouldMatchOnlyGivenString(t *testing.T) {
	r := text.NewReader([]byte("bar"), true)
	_, res := terminal.Word("foo", "TEST", 5).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
}
