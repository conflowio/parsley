package parser_test

import (
	"testing"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/parser"
	"github.com/opsidian/parsec/reader"
	"github.com/stretchr/testify/assert"
)

func TestMemoizeShouldIncDecCallCount(t *testing.T) {
	r := reader.New([]byte("aa"), true)
	r.ReadRune()
	c := parser.NewContext()
	assert.Equal(t, c.GetCalls("p1", 1), 0)

	p := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		assert.Equal(t, c.GetCalls("p1", 1), 1)
		return nil
	})
	parser.Memoize("p1", p)(c, r)

	assert.Equal(t, c.GetCalls("p1", 1), 0)
}

func TestMemoizeShouldReturnParserResult(t *testing.T) {
	r := reader.New([]byte("a"), true)
	c := parser.NewContext()

	node := ast.NewTerminalNode("a", reader.NewPosition(0, 1, 2), "a")
	expected := parser.NewResult(node, r).AsList()

	p := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		return expected
	})
	actual := parser.Memoize("p1", p)(c, r)
	assert.Equal(t, expected, actual)
}

func TestMemoizeShouldRememberResult(t *testing.T) {
	r := reader.New([]byte("a"), true)
	c := parser.NewContext()

	node := ast.NewTerminalNode("t", reader.NewPosition(0, 1, 2), nil)
	expected := parser.NewResult(node, r).AsList()

	called := false
	p := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		called = true
		return expected
	})
	parser.Memoize("p1", p)(c, r)

	called = false

	actual := parser.Memoize("p1", p)(c, r)

	assert.False(t, called, "The parser function should not have been called")
	assert.Equal(t, expected, actual)

	r.ReadRune()

	parser.Memoize("p1", p)(c, r)
	assert.True(t, called, "The parser function should have been called for the new position")
}

// The parser calls should be curtailed if the parser was already called more times than the remaining token count plus one
func TestMemoizeShouldCurtailResult(t *testing.T) {
	r := reader.New([]byte("a"), true)
	c := parser.NewContext()
	c.Push("p1", 0)
	c.Push("p1", 0)
	c.Push("p2", 0)
	called := false
	p := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		called = true
		return nil
	})
	expected := parser.NewResults(map[int]int{0: 3, 1: 1})
	results := parser.Memoize("p1", p)(c, r)
	assert.False(t, called, "The call tree should have been curtailed")
	assert.Equal(t, expected, results)
}

func TestOrShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()
	assert.Panics(t, func() { parser.Or("o1")(c, r) })
}

func TestOrShouldHandleOnlyOneParser(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()

	var r1 parser.Result

	p1 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r)
		return r1.AsList()
	})

	results := parser.Or("o1", p1)(c, r)
	assert.Equal(t, parser.NewResults(nil, r1), results)
}

func TestOrShouldMergeResults(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()

	var r1, r2 parser.Result

	p1 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r)
		return r1.AsList()
	})

	p2 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		pos := r.Cursor()
		ch1, _, _ := r.ReadRune()
		ch2, _, _ := r.ReadRune()
		r2 = parser.NewResult(ast.NewTerminalNode("STRING", pos, string([]rune{ch1, ch2})), r)
		return r2.AsList()
	})

	p3 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		return c.NewCurtailedResults(r.Cursor().Pos())
	})

	p4 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		return nil
	})

	p5 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		return parser.NewResults(nil)
	})

	results := parser.Or("o1", p1, p2, p3, p4, p5)(c, r)
	assert.EqualValues(t, parser.NewResults(map[int]int{0: 1}, r1, r2), results)

	assert.Equal(t, 6, c.GetSumCallCount())
}

func TestOrMayReturnEmptyResult(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()

	p1 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		return nil
	})

	p2 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		return nil
	})

	results := parser.Or("o1", p1, p2)(c, r)
	assert.Equal(t, parser.NewResults(nil), results)
}

func TestOrShouldCloneReadersForAllParsers(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()

	p1 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		assert.Equal(t, reader.NewPosition(0, 1, 1), r.Cursor())
		r.ReadRune()
		return nil
	})

	p2 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		assert.Equal(t, reader.NewPosition(0, 1, 1), r.Cursor())
		r.ReadRune()
		return nil
	})

	parser.Or("o1", p1, p2)(c, r)
	assert.Equal(t, 3, c.GetSumCallCount())
}

func TestAndShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()
	assert.Panics(t, func() { parser.And("o1", nil)(c, r) })
}

func TestAndShouldHandleOnlyOneParser(t *testing.T) {
	r := reader.New([]byte("ab"), true)
	c := parser.NewContext()

	var r1 parser.Result

	p1 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r)
		return r1.AsList()
	})

	nodeBuilder := func(nodes []ast.Node) ast.Node { return nodes[0] }

	results := parser.And("o1", nodeBuilder, p1)(c, r)
	assert.Equal(t, parser.NewResults(nil, r1), results)
}

func TestAndShouldCombineParserResults(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	c := parser.NewContext()

	// This function simulates a complex parser which returns two different results
	p := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		pos := r.Cursor()
		r1 := r.Clone()
		ch1, _, _ := r1.ReadRune()
		res1 := parser.NewResult(ast.NewTerminalNode("CHAR", pos, string(ch1)), r1)
		r2 := r1.Clone()
		ch2, _, _ := r2.ReadRune()
		res2 := parser.NewResult(ast.NewTerminalNode("STRING", pos, string([]rune{ch1, ch2})), r2)
		return parser.NewResults(nil, res1, res2)
	})

	nodeBuilder := func(nodes []ast.Node) ast.Node {
		var res string
		for _, node := range nodes {
			val, _ := node.Value()
			res += val.(string)
		}
		return ast.NewTerminalNode("STRING", nodes[0].Pos(), res)
	}

	results := parser.And("o1", nodeBuilder, p, p)(c, r)
	assert.EqualValues(t, parser.NewResults(
		nil,
		parser.NewResult(ast.NewTerminalNode("STRING", reader.NewPosition(0, 1, 1), "ab"), r.WithCursor(2, 1, 3)),
		parser.NewResult(ast.NewTerminalNode("STRING", reader.NewPosition(0, 1, 1), "abc"), r.WithCursor(3, 1, 4)),
		parser.NewResult(ast.NewTerminalNode("STRING", reader.NewPosition(0, 1, 1), "abc"), r.WithCursor(3, 1, 4)),
		parser.NewResult(ast.NewTerminalNode("STRING", reader.NewPosition(0, 1, 1), "abcd"), r.WithCursor(4, 1, 5)),
	), results)

	assert.EqualValues(t, 4, c.GetSumCallCount())
}

func TestAndShouldHandleNilResults(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	c := parser.NewContext()

	p1 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		return parser.NewResults(nil, parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r))
	})

	p2 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		return nil
	})

	results := parser.And("o1", nil, p1, p2)(c, r)
	assert.EqualValues(t, parser.NewResults(nil), results)
}

func TestAndShouldMergeCurtailReasonsIfEmptyResult(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	c := parser.NewContext()

	p1 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		return parser.NewResults(map[int]int{0: 1, 1: 2}, parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r))
	})

	p2 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		return parser.NewResults(map[int]int{0: 2, 1: 1})
	})

	results := parser.And("o1", nil, p1, p2)(c, r)
	assert.EqualValues(t, parser.NewResults(map[int]int{0: 2, 1: 2}), results)
}

func TestAndShouldStopIfEOFReached(t *testing.T) {
	r := reader.New([]byte("abcd"), true)
	c := parser.NewContext()

	p1 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		return parser.NewResults(nil, parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r))
	})

	p2 := parser.Func(func(c *parser.Context, r *reader.Reader) (results *parser.Results) {
		pos := r.Cursor()
		r1 := r.Clone()
		ch, _, _ := r1.ReadRune()
		return parser.NewResults(nil,
			parser.NewResult(ast.NewTerminalNode(reader.EOF, pos, nil), r),
			parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r1),
		)
	})

	nodeBuilder := func(nodes []ast.Node) ast.Node {
		return nodes[0]
	}

	results := parser.And("o1", nodeBuilder, p1, p2)(c, r)
	assert.EqualValues(t, parser.NewResults(
		nil,
		parser.NewResult(ast.NewTerminalNode("CHAR", reader.NewPosition(0, 1, 1), 'a'), r.WithCursor(1, 1, 2)),
	), results)
}
