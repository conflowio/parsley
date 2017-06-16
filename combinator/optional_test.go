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

func TestOptionalShouldReturnParserResult(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	res := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), test.NewReader(1, 1, false, true))

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return data.NewIntSet(1), res.AsSet()
	})

	cp, rs := combinator.Optional(p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, data.NewIntSet(1), cp)
	assert.Equal(t, parser.NewResultSet(res), rs)
}

func TestOptionalShouldReturnEmptyResultIfParserFailed(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return data.NewIntSet(1), nil
	})

	cp, rs := combinator.Optional(p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, data.NewIntSet(1), cp)
	assert.Equal(t, parser.NewResultSet(parser.NewResult(nil, r.Clone())), rs)
}
