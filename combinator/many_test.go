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
	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		defer func() { pi++ }()
		if pi < len(pResults) {
			return parser.NoCurtailingParsers(), pResults[pi]
		} else {
			return parser.NoCurtailingParsers(), nil
		}
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		var res string
		for _, node := range nodes {
			val, _ := node.Value()
			res += val.(string)
		}
		return ast.NewTerminalNode("STR", nodes[0].Pos(), res)
	})

	_, rs := combinator.Many(nodeBuilder, p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, parser.NewResultSet(
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "ac"), test.NewReader(3, 1, false, true)),
		parser.NewResult(ast.NewTerminalNode("STR", test.NewPosition(1), "bd"), test.NewReader(4, 1, false, true)),
	), rs)
}

func TestManyShouldHandleNilResults(t *testing.T) {
	r := test.NewReader(0, 1, false, false)

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		return parser.NoCurtailingParsers(), nil
	})

	cp, rs := combinator.Many(nil, p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
}

func TestManyShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
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
