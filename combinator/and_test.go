package combinator_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
)

func TestAndShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	ctx := parser.EmptyLeftRecCtx()
	assert.Panics(t, func() { combinator.And(nil).Parse(ctx, r) })
}

func TestAndShouldHandleOnlyOneParser(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	ctx := parser.EmptyLeftRecCtx()

	var r1 parser.Result
	curtailingParsers := data.NewIntSet(1)

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r)
		return parser.NewParserResult(curtailingParsers, r1)
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node { return nodes[0] })

	results := combinator.And(nodeBuilder, p1).Parse(ctx, r)
	assert.Equal(t, parser.NewParserResult(curtailingParsers, r1), results)
}

func TestAndShouldCombineParserResults(t *testing.T) {
	parser.Stat.Reset()
	r := reader.New([]byte("abcd"), true)
	ctx := parser.EmptyLeftRecCtx()

	// This function simulates a complex parser which returns two different results
	p := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		r1 := r.Clone()
		ch1, _, _ := r1.ReadRune()
		res1 := parser.NewResult(ast.NewTerminalNode("CHAR", pos, string(ch1)), r1)
		r2 := r1.Clone()
		ch2, _, _ := r2.ReadRune()
		res2 := parser.NewResult(ast.NewTerminalNode("STRING", pos, string([]rune{ch1, ch2})), r2)
		return parser.NewParserResult(parser.NoCurtailingParsers(), res1, res2)
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		var res string
		for _, node := range nodes {
			val, _ := node.Value()
			res += val.(string)
		}
		return ast.NewTerminalNode("STRING", nodes[0].Pos(), res)
	})

	results := combinator.And(nodeBuilder, p, p).Parse(ctx, r)
	assert.EqualValues(t, parser.NewParserResult(
		parser.NoCurtailingParsers(),
		parser.NewResult(ast.NewTerminalNode("STRING", reader.NewPosition(0, 1, 1), "ab"), r.WithCursor(2, 1, 3)),
		parser.NewResult(ast.NewTerminalNode("STRING", reader.NewPosition(0, 1, 1), "abc"), r.WithCursor(3, 1, 4)),
		parser.NewResult(ast.NewTerminalNode("STRING", reader.NewPosition(0, 1, 1), "abcd"), r.WithCursor(4, 1, 5)),
	), results)

	assert.EqualValues(t, 3, parser.Stat.GetSumCallCount())
}

func TestAndShouldHandleNilResults(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	ctx := parser.EmptyLeftRecCtx()

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		return parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r))
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return nil
	})

	results := combinator.And(nil, p1, p2).Parse(ctx, r)
	assert.EqualValues(t, parser.NewParserResult(parser.NoCurtailingParsers()), results)
}

func TestAndShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	ctx := parser.EmptyLeftRecCtx()

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return parser.NewParserResult(data.NewIntSet(0, 1), parser.NewResult(nil, r))
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return parser.NewParserResult(data.NewIntSet(1, 2))
	})

	results := combinator.And(nil, p1, p2).Parse(ctx, r)
	assert.EqualValues(t, parser.NewParserResult(data.NewIntSet(0, 1, 2)), results)
}

func TestAndShouldStopIfEOFReached(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	ctx := parser.EmptyLeftRecCtx()

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		return parser.NewParserResult(parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r))
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		r1 := r.Clone()
		ch, _, _ := r1.ReadRune()
		return parser.NewParserResult(parser.NoCurtailingParsers(),
			parser.NewResult(ast.NewTerminalNode(ast.EOF, pos, nil), r),
			parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r1),
		)
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		return nodes[0]
	})

	results := combinator.And(nodeBuilder, p1, p2).Parse(ctx, r)
	assert.EqualValues(t, parser.NewParserResult(
		parser.NoCurtailingParsers(),
		parser.NewResult(ast.NewTerminalNode("CHAR", reader.NewPosition(0, 1, 1), 'a'), r.WithCursor(1, 1, 2)),
	), results)
}
