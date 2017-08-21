// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	"fmt"
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/test"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Using the Empty parser you can match expressions optionally.
// Note: there is an Optional combinator for this purpose.
func ExampleEmpty() {
	concat := ast.InterpreterFunc(func(ctx interface{}, nodes []ast.Node) (interface{}, reader.Error) {
		var res string
		for _, node := range nodes {
			if node != nil {
				val, _ := node.Value(ctx)
				res += string(val.(rune))
			}
		}
		return res, nil
	})

	p := combinator.Seq(builder.All("ABC", concat),
		terminal.Rune('a', "a"),
		combinator.Choice("b or nothing", terminal.Rune('b', "b"), parser.Empty()),
		terminal.Rune('c', "c"),
	)
	s := parsley.NewSentence(p)
	value, _, _ := s.Evaluate(text.NewReader([]byte("ac"), true), nil)
	fmt.Printf("%T %v\n", value, value)
	// Output: string ac
}

// Using the End parser you can make sure you fully match the input
func ExampleEnd() {
	s := combinator.Seq(builder.Select(0), terminal.Float(), parser.End())
	_, rs, _ := s.Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), text.NewReader([]byte("1.23"), true))
	value, _ := rs[0].Node().Value(nil)
	fmt.Printf("%T %v\n", value, value)
	// Output: float64 1.23
}

func assertCursor(t *testing.T, pos int, r reader.Reader) {
	assert.Equal(t, pos, r.Cursor().Pos())
}

func TestEmptyWillAlwaysReturnWithResult(t *testing.T) {
	r := test.NewReader(1, 1, false, false)
	_, res, err := parser.Empty()(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, parser.NewResult(nil, r).AsSet(), res)
	assertCursor(t, 1, r)
	assert.Nil(t, err)
}

func TestEndShouldMatchEOF(t *testing.T) {
	r := test.NewReader(1, 1, true, false)
	_, res, err := parser.End()(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode(ast.EOF, test.NewPosition(1), nil)
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assertCursor(t, 1, r)
	assert.Nil(t, err)
}

func TestEndShouldNotMatchNotEOF(t *testing.T) {
	r := test.NewReader(1, 1, false, false)
	_, res, err := parser.End()(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	require.NotNil(t, err)
	assert.Equal(t, test.NewPosition(1), err.Pos())
}
