package combinator_test

import (
	"fmt"
	"testing"

	"github.com/opsidian/parsley"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/test"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Let's define a parser which accepts integer or float numbers.
// The parser would return only the first match so in this case we have to put the float parser first.
func ExampleChoice() {
	number := combinator.Choice("number",
		terminal.Float(),
		terminal.Integer(),
	)
	s := combinator.Seq(builder.Select(0), number, parser.End())
	value, _ := parsley.EvaluateText([]byte("1.23"), true, s, nil)
	fmt.Printf("%T %v\n", value, value)
	// Output: float64 1.23
}

func TestChoiceShouldPanicIfNoParserWasGiven(t *testing.T) {
	r := test.NewReader(0, 2, false, false)
	assert.Panics(t, func() { combinator.Choice("test").Parse(parser.EmptyLeftRecCtx(), r) })
}

func TestChoiceShouldHandleOnlyOneParser(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	expectedCP := data.NewIntSet(1)
	expectedRS := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'x'), r).AsSet()

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return expectedCP, expectedRS, nil
	})

	cp, rs, err := combinator.Choice("test", p1).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, expectedCP, cp)
	assert.Equal(t, expectedRS, rs)
	assert.Nil(t, err)
}

func TestChoiceShouldMergeCurtailingParsers(t *testing.T) {
	parser.Stat.Reset()
	r := test.NewReader(0, 2, false, false)

	var res parser.Result

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(test.NewPosition(1), "ERR1")
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return data.NewIntSet(1), parser.NewResultSet(), parser.NewError(test.NewPosition(2), "ERR2")
	})

	p3 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		res = parser.NewResult(ast.NewTerminalNode("STRING", test.NewPosition(0), "TEST"), test.NewReader(1, 1, false, true))
		return data.NewIntSet(2), res.AsSet(), nil
	})

	cp, rs, err := combinator.Choice("test", p1, p2, p3).Parse(parser.EmptyLeftRecCtx(), r)
	expectedCP := data.NewIntSet(1, 2)
	expectedRS := res.AsSet()
	assert.EqualValues(t, expectedCP, cp)
	assert.EqualValues(t, expectedRS, rs)
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(2), err.Pos())
	assert.Equal(t, "ERR2 at Pos{2}", err.Error())

	assert.Equal(t, 3, parser.Stat.GetSumCallCount())
}

func TestChoiceShouldStopAtChoiceResult(t *testing.T) {
	parser.Stat.Reset()
	r := test.NewReader(0, 2, false, false)

	var res1, res2 parser.Result

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		res1 = parser.NewResult(ast.NewTerminalNode("STRING", test.NewPosition(0), "TEST"), test.NewReader(1, 1, false, true))
		return parser.NoCurtailingParsers(), res1.AsSet(), nil
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		res2 = parser.NewResult(ast.NewTerminalNode("STRING", test.NewPosition(0), "TEST2"), test.NewReader(1, 1, false, true))
		return parser.NoCurtailingParsers(), res2.AsSet(), nil
	})

	_, rs, err := combinator.Choice("test", p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.EqualValues(t, res1.AsSet(), rs)
	assert.Nil(t, err)

	assert.Equal(t, 1, parser.Stat.GetSumCallCount())
}

func TestChoiceMayReturnEmptyResult(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(test.NewPosition(2), "ERR1")
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(test.NewPosition(1), "ERR2")
	})

	cp, rs, err := combinator.Choice("test", p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(2), err.Pos())
	assert.Equal(t, "ERR1 at Pos{2}", err.Error())
}

func TestChoiceShouldReturnCustomErrorIfNoParserAdvanced(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	p1 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(test.NewPosition(0), "ERR1")
	})

	p2 := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		return parser.NoCurtailingParsers(), nil, parser.NewError(test.NewPosition(0), "ERR2")
	})

	cp, rs, err := combinator.Choice("test", p1, p2).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NoCurtailingParsers(), cp)
	assert.Empty(t, rs)
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(0), err.Pos())
	assert.Equal(t, "was expecting test at Pos{0}", err.Error())
}
