package combinator_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestMaybeShouldReturnParserResultAndEmptyResult(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	r1 := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), test.NewReader(1, 1, false, true))

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return data.NewIntSet(1), r1.AsSet()
	})
	r2 := parser.NewResult(ast.NewTerminalNode(ast.EMPTY, r.Cursor(), nil), r.Clone())

	cp, rs := combinator.Maybe(p1).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, data.NewIntSet(1), cp)
	assert.Equal(t, parser.NewResultSet(r2, r1), rs)
}
