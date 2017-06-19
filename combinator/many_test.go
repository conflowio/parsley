package combinator_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManyShouldCombineParserResults(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	pResults := []parser.ResultSet{
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "a"), test.NewReader(1, 1, false, true)),
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "b"), test.NewReader(2, 1, false, true)),
		),
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "c"), test.NewReader(3, 1, false, true)),
		),
		nil,
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "d"), test.NewReader(4, 1, false, true)),
		),
		nil,
	}

	pi := 0
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		defer func() { pi++ }()
		if pi < len(pResults) {
			return parser.NoCurtailingParsers(), pResults[pi], nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		var res string
		for _, node := range nodes {
			val, _ := node.Value()
			res += val.(string)
		}
		first := nodes[0].(ast.TerminalNode)
		return ast.NewTerminalNode("STR", first.Pos(), res)
	})

	_, rs, err := combinator.Many(nodeBuilder, p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, parser.NewResultSet(
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(3, 1, false, true)),
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "bd"), test.NewReader(4, 1, false, true)),
	), rs)
	assert.Nil(t, err)

	pi = 0 // Reset parser call index
	_, rs, err = combinator.Many1(nodeBuilder, p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, parser.NewResultSet(
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(3, 1, false, true)),
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "bd"), test.NewReader(4, 1, false, true)),
	), rs)
	assert.Nil(t, err)
}

func TestMany1ShouldReturnNoResultIfNoMatch(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(test.NewPosition(1), "ERR1")
	})

	cp, rs, err := combinator.Many1(builder.Nil(), p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(1), err.Pos())
}

func TestManyShouldReturnEmptyResultIfNoMatch(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(test.NewPosition(1), "ERR1")
	})

	cp, rs, err := combinator.Many(builder.Nil(), p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Equal(t, rs, parser.NewResult(nil, r).AsSet())
	assert.Nil(t, err)
}

func TestManyShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	pi := 0
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		defer func() { pi++ }()
		if pi == 0 {
			return data.NewIntSet(0, 1), parser.NewResult(nil, r).AsSet(), nil
		} else {
			return data.NewIntSet(1, 2), nil, parser.NewError(test.NewPosition(1), "ERR1")
		}
	})

	cp, _, _ := combinator.Many(builder.Nil(), p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)

	pi = 0 // Reset parser call index
	cp, _, _ = combinator.Many1(builder.Nil(), p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
}
