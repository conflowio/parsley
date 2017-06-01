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

func TestManySepByShouldCombineParserResults(t *testing.T) {
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
		nil,
	}

	pi := 0
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		defer func() { pi++ }()
		if pi < len(pResults) {
			return parser.NoCurtailingParsers(), pResults[pi]
		} else {
			return parser.NoCurtailingParsers(), nil
		}
	})

	sepi := 0
	sep := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		defer func() { sepi++ }()
		if sepi < len(pResults) {
			return parser.NoCurtailingParsers(), sepResults[sepi]
		} else {
			return parser.NoCurtailingParsers(), nil
		}
	})

	interpreter := ast.InterpreterFunc(func(values []interface{}) (interface{}, error) {
		res := ""
		for _, value := range values {
			res += value.(string) + "|"
		}
		return res, nil
	})

	_, rs := combinator.ManySepBy("TEST", "X", h, p, sep, interpreter).Parse(parser.EmptyLeftRecCtx(), r)
	val1, _ := rs[0].Node().Value()
	val2, _ := rs[1].Node().Value()
	assert.Equal(t, "a|,|c|", val1)
	assert.Equal(t, "b|,|d|", val2)
}

func TestManySepByShouldHandleNilResults(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})

	cp, rs := combinator.Many(nil, p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
}

func TestManySepByShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	pi := 0
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		defer func() { pi++ }()
		if pi == 0 {
			return data.NewIntSet(0, 1), parser.NewResult(ast.NewTerminalNode(ast.EMPTY, test.NewPosition(0), nil), r).AsSet()
		} else {
			return data.NewIntSet(1, 2), nil
		}
	})

	cp, _ := combinator.Many(ast.NilBuilder(), p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, data.NewIntSet(0, 1, 2), cp)
}
