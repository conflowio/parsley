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

func TestMaybeShouldReturnParserResultAndEmptyResult(t *testing.T) {
	r := reader.New([]byte("ab"), true)

	var r1 parser.Result

	p1 := parser.Func(func(ctx data.IntMap, r *reader.Reader) (results *parser.ParserResult) {
		pos := r.Cursor()
		ch, _, _ := r.ReadRune()
		r1 = parser.NewResult(ast.NewTerminalNode("CHAR", pos, ch), r)
		return parser.NewParserResult(parser.NoCurtailingParsers(), r1)
	})
	r2 := parser.NewResult(ast.NewTerminalNode(ast.EMPTY, r.Cursor(), nil), r)

	results := combinator.Maybe(p1).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NewParserResult(parser.NoCurtailingParsers(), r2, r1), results)
}
