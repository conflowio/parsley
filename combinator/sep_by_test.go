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
	"github.com/stretchr/testify/require"
)

func TestSepByShouldCombineParserResults(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	h := parser.NewHistory()

	pResults := []parser.ResultSet{
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(0), "a"), test.NewReader(1, 1, false, true)),
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "b"), test.NewReader(2, 1, false, true)),
		),
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(4), "c"), test.NewReader(5, 1, false, true)),
		),
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(5), "d"), test.NewReader(6, 1, false, true)),
		),
	}

	sepResults := []parser.ResultSet{
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), ","), test.NewReader(3, 1, false, true)),
		),
		nil,
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(3), ","), test.NewReader(4, 1, false, true)),
		),
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

	sepi := 0
	sep := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		defer func() { sepi++ }()
		if sepi < len(sepResults) {
			return parser.NoCurtailingParsers(), sepResults[sepi], nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	interpreter := ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
		res := ""
		for _, node := range nodes {
			value, _ := node.Value()
			res += value.(string) + "|"
		}
		return res, nil
	})

	_, rs, err := combinator.SepBy1("TEST", h, p, sep, interpreter).Parse(parser.EmptyLeftRecCtx(), r)
	require.Len(t, rs, 2)
	val0, _ := rs[0].Node().Value()
	val1, _ := rs[1].Node().Value()
	assert.Equal(t, "a|,|c|", val0)
	assert.Equal(t, "b|,|d|", val1)
	assert.Nil(t, err)
}

func TestSepByShouldNotFlattenNonTerminals(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	h := parser.NewHistory()

	interpreter1 := ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
		res := ""
		for _, node := range nodes {
			value, _ := node.Value()
			res += value.(string) + "&"
		}
		return res, nil
	})

	pResults := []parser.ResultSet{
		parser.NewResultSet(
			parser.NewResult(
				ast.NewNonTerminalNode("NT", []ast.Node{
					ast.NewTerminalNode("STR", test.NewPosition(0), "a"),
					ast.NewTerminalNode("STR", test.NewPosition(1), "b"),
				}, interpreter1),
				test.NewReader(2, 1, false, true),
			),
		),
		parser.NewResultSet(
			parser.NewResult(
				ast.NewNonTerminalNode("NT", []ast.Node{
					ast.NewTerminalNode("STR", test.NewPosition(3), "c"),
					ast.NewTerminalNode("STR", test.NewPosition(4), "d"),
				}, interpreter1),
				test.NewReader(5, 1, false, true),
			),
		),
	}

	sepResults := []parser.ResultSet{
		parser.NewResultSet(
			parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(2), ","), test.NewReader(3, 1, false, true)),
		),
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

	sepi := 0
	sep := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		defer func() { sepi++ }()
		if sepi < len(sepResults) {
			return parser.NoCurtailingParsers(), sepResults[sepi], nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	interpreter2 := ast.InterpreterFunc(func(nodes []ast.Node) (interface{}, error) {
		res := ""
		for _, node := range nodes {
			value, _ := node.Value()
			res += value.(string) + "|"
		}
		return res, nil
	})

	_, rs, _ := combinator.SepBy1("TEST", h, p, sep, interpreter2).Parse(parser.EmptyLeftRecCtx(), r)
	require.Len(t, rs, 1)
	val0, _ := rs[0].Node().Value()
	assert.Equal(t, "a&b&|,|c&d&|", val0)
}

func TestSepByShouldReturnEmptyResultIfNoMatch(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	h := parser.NewHistory()

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(test.NewPosition(1), "TEST1")
	})

	cp, rs, err := combinator.SepBy("TEST", h, p, p, nil).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Equal(t, parser.NewResult(ast.NewNonTerminalNode("TEST", nil, nil), r).AsSet(), rs)
	assert.Nil(t, err)
}

func TestSepBy1ShouldReturnNilIfNoResult(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	h := parser.NewHistory()

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(test.NewPosition(1), "TEST1")
	})

	cp, rs, err := combinator.SepBy1("TEST", h, p, p, nil).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(1), err.Pos())
}

func TestSepByShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
	r := test.NewReader(0, 1, false, false)
	h := parser.NewHistory()

	pi := 0
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		defer func() { pi++ }()
		if pi < 1 {
			return data.NewIntSet(0, 1), parser.NewResult(nil, r).AsSet(), nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	sepi := 0
	sep := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		defer func() { sepi++ }()
		if sepi < 1 {
			return data.NewIntSet(1, 2), nil, nil
		} else {
			return parser.NoCurtailingParsers(), nil, nil
		}
	})

	cp, _, _ := combinator.SepBy1("TEST", h, p, sep, nil).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
}
