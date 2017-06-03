package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func assertCursor(t *testing.T, pos int, r reader.Reader) {
	assert.Equal(t, pos, r.Cursor().Pos())
}

func TestEmptyWillAlwaysReturnWithResult(t *testing.T) {
	r := test.NewReader(1, 1, false, false)
	_, res := parser.Empty()(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NewResult(nil, r).AsSet(), res)
	assertCursor(t, 1, r)
}

func TestEndShouldMatchEOF(t *testing.T) {
	r := test.NewReader(1, 1, true, false)
	_, res := parser.End()(parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode(ast.EOF, test.NewPosition(1), nil)
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assertCursor(t, 1, r)
}

func TestEndShouldNotMatchNotEOF(t *testing.T) {
	r := test.NewReader(1, 1, false, false)
	_, res := parser.End()(parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
}
