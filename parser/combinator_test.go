package parser_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/stretchr/testify/assert"
)

func TestMemoizeShouldIncreaseLeftRecContext(t *testing.T) {
	r := reader.New([]byte("aa"), true)
	r.ReadRune()
	c := parser.NewContext()
	leftRecContext := data.NewIntMap()
	parserIndex := c.GetParserIndex("p1")
	assert.Equal(t, leftRecContext.Get(parserIndex), 0)

	p := parser.Func(func(leftRecContext data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		assert.Equal(t, leftRecContext.Get(parserIndex), 1)
		return nil
	})
	parser.Memoize("p1", c, p).Parse(leftRecContext, r)
}

func TestMemoizeShouldReturnParserResult(t *testing.T) {
	r := reader.New([]byte("a"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	node := ast.NewTerminalNode("a", reader.NewPosition(0, 1, 2), "a")
	expected := parser.NewParserResult(data.NewIntSet(), parser.NewResult(node, r))

	p := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return expected
	})
	actual := parser.Memoize("p1", c, p).Parse(ctx, r)
	assert.Equal(t, expected, actual)
}

func TestMemoizeShouldRememberResult(t *testing.T) {
	r := reader.New([]byte("a"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	expected := parser.NewParserResult(data.NewIntSet(), parser.NewResult(node, r))

	called := false
	p := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		called = true
		return expected
	})
	parser.Memoize("p1", c, p).Parse(ctx, r)

	called = false

	actual := parser.Memoize("p1", c, p).Parse(ctx, r)

	assert.False(t, called, "The parser function should not have been called")
	assert.Equal(t, expected, actual)

	r.ReadRune()

	parser.Memoize("p1", c, p).Parse(ctx, r)
	assert.True(t, called, "The parser function should have been called for the new position")
}

// The parser calls should be curtailed if the parser was already called more times than the remaining token count plus one
func TestMemoizeShouldCurtailResult(t *testing.T) {
	r := reader.New([]byte("a"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap().
		Inc(c.GetParserIndex("p1")).
		Inc(c.GetParserIndex("p1")).
		Inc(c.GetParserIndex("p1")).
		Inc(c.GetParserIndex("p2"))

	called := false
	p := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		called = true
		return nil
	})
	expected := parser.NewParserResult(data.NewIntSet().Insert(c.GetParserIndex("p1")))
	results := parser.Memoize("p1", c, p).Parse(ctx, r)
	assert.False(t, called, "The call tree should have been curtailed")
	assert.Equal(t, expected, results)
}

func TestOrShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()
	assert.Panics(t, func() { parser.Or("o1", c).Parse(ctx, r) })
}

func TestOrShouldHandleOnlyOneParser(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	var r1 parser.Result

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r)
		return parser.NewParserResult(data.NewIntSet(), r1)
	})

	results := parser.Or("o1", c, p1).Parse(ctx, r)
	assert.Equal(t, parser.NewParserResult(data.NewIntSet(), r1), results)
}

func TestOrShouldMergeResults(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	var r1, r2 parser.Result

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r)
		return parser.NewParserResult(data.NewIntSet(), r1)
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch1, _, _ := r.ReadRune()
		ch2, _, _ := r.ReadRune()
		r2 = parser.NewResult(ast.NewTerminalNode("STRING", pos, string([]rune{ch1, ch2})), r)
		return parser.NewParserResult(data.NewIntSet().Insert(1), r2)
	})

	p3 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return parser.NewParserResult(data.NewIntSet().Insert(2))
	})

	p4 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return nil
	})

	results := parser.Or("o1", c, p1, p2, p3, p4).Parse(ctx, r)
	curtailingParsers := data.NewIntSet().Insert(1).Insert(2)
	assert.EqualValues(t, parser.NewParserResult(curtailingParsers, r1, r2), results)

	assert.Equal(t, 5, c.GetSumCallCount())
}

func TestOrMayReturnEmptyResult(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return nil
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return nil
	})

	results := parser.Or("o1", c, p1, p2).Parse(ctx, r)
	assert.Equal(t, parser.NewParserResult(data.NewIntSet()), results)
}

func TestOrShouldCloneReadersForAllParsers(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		assert.Equal(t, reader.NewPosition(0, 1, 1), r.Cursor())
		r.ReadRune()
		return nil
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		assert.Equal(t, reader.NewPosition(0, 1, 1), r.Cursor())
		r.ReadRune()
		return nil
	})

	parser.Or("o1", c, p1, p2).Parse(ctx, r)
	assert.Equal(t, 3, c.GetSumCallCount())
}

func TestAndShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()
	assert.Panics(t, func() { parser.And("o1", c, nil).Parse(ctx, r) })
}

func TestAndShouldHandleOnlyOneParser(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	var r1 parser.Result
	curtailingParsers := data.NewIntSet().Insert(1)

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r)
		return parser.NewParserResult(curtailingParsers, r1)
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node { return nodes[0] })

	results := parser.And("o1", c, nodeBuilder, p1).Parse(ctx, r)
	assert.Equal(t, parser.NewParserResult(curtailingParsers, r1), results)
}

func TestAndShouldCombineParserResults(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	// This function simulates a complex parser which returns two different results
	p := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		r1 := r.Clone()
		ch1, _, _ := r1.ReadRune()
		res1 := parser.NewResult(ast.NewTerminalNode("CHAR", pos, string(ch1)), r1)
		r2 := r1.Clone()
		ch2, _, _ := r2.ReadRune()
		res2 := parser.NewResult(ast.NewTerminalNode("STRING", pos, string([]rune{ch1, ch2})), r2)
		return parser.NewParserResult(data.NewIntSet(), res1, res2)
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		var res string
		for _, node := range nodes {
			val, _ := node.Value()
			res += val.(string)
		}
		return ast.NewTerminalNode("STRING", nodes[0].Pos(), res)
	})

	results := parser.And("o1", c, nodeBuilder, p, p).Parse(ctx, r)
	assert.EqualValues(t, parser.NewParserResult(
		data.NewIntSet(),
		parser.NewResult(ast.NewTerminalNode("STRING", reader.NewPosition(0, 1, 1), "ab"), r.WithCursor(2, 1, 3)),
		parser.NewResult(ast.NewTerminalNode("STRING", reader.NewPosition(0, 1, 1), "abc"), r.WithCursor(3, 1, 4)),
		parser.NewResult(ast.NewTerminalNode("STRING", reader.NewPosition(0, 1, 1), "abcd"), r.WithCursor(4, 1, 5)),
	), results)

	assert.EqualValues(t, 4, c.GetSumCallCount())
}

func TestAndShouldHandleNilResults(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		return parser.NewParserResult(data.NewIntSet(), parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r))
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return nil
	})

	results := parser.And("o1", c, nil, p1, p2).Parse(ctx, r)
	assert.EqualValues(t, parser.NewParserResult(data.NewIntSet()), results)
}

func TestAndShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return parser.NewParserResult(data.NewIntSet(0, 1), parser.NewResult(nil, r))
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		return parser.NewParserResult(data.NewIntSet(1, 2))
	})

	results := parser.And("o1", c, nil, p1, p2).Parse(ctx, r)
	assert.EqualValues(t, parser.NewParserResult(data.NewIntSet(0, 1, 2)), results)
}

func TestAndShouldStopIfEOFReached(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	c := parser.NewContext()
	ctx := data.NewIntMap()

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		return parser.NewParserResult(data.NewIntSet(), parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r))
	})

	p2 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		r1 := r.Clone()
		ch, _, _ := r1.ReadRune()
		return parser.NewParserResult(data.NewIntSet(),
			parser.NewResult(ast.NewTerminalNode(reader.EOF, pos, nil), r),
			parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r1),
		)
	})

	nodeBuilder := ast.NodeBuilderFunc(func(nodes []ast.Node) ast.Node {
		return nodes[0]
	})

	results := parser.And("o1", c, nodeBuilder, p1, p2).Parse(ctx, r)
	assert.EqualValues(t, parser.NewParserResult(
		data.NewIntSet(),
		parser.NewResult(ast.NewTerminalNode("CHAR", reader.NewPosition(0, 1, 1), 'a'), r.WithCursor(1, 1, 2)),
	), results)
}
