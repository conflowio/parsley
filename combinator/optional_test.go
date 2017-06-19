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

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return data.NewIntSet(1), res.AsSet(), nil
	})

	cp, rs, err := combinator.Optional(p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, data.NewIntSet(1), cp)
	assert.Equal(t, parser.NewResultSet(res), rs)
	assert.Nil(t, err)
}

func TestOptionalShouldReturnEmptyResultIfParserFailed(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return data.NewIntSet(1), nil, nil
	})

	cp, rs, err := combinator.Optional(p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, data.NewIntSet(1), cp)
	assert.Equal(t, parser.NewResultSet(parser.NewResult(nil, r.Clone())), rs)
	assert.Nil(t, err)
}
