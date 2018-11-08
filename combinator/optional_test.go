// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator_test

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Let's define a parser which accepts "a", an optional "b" and a "c" character.
// The optional parser will result in a nil node so in the interpreter we have to handle that.
func ExampleOptional() {
	concat := ast.InterpreterFunc(func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
		var res string
		for _, node := range node.Children() {
			if node != nil {
				val, _ := node.Value(userCtx)
				if val != nil {
					res += string(val.(rune))
				}
			}
		}
		return res, nil
	})

	p := combinator.SeqOf(
		terminal.Rune('a'),
		combinator.Optional(terminal.Rune('b')),
		terminal.Rune('c'),
	).Bind(concat)
	r := text.NewReader(text.NewFile("example.file", []byte("ac")))
	ctx := parsley.NewContext(parsley.NewFileSet(), r)
	value, _ := parsley.Evaluate(ctx, combinator.Sentence(p), nil)
	fmt.Printf("%T %v\n", value, value)
	// Output: string ac
}

//
// func TestOptionalShouldReturnParserResult(t *testing.T) {
// 	r := test.NewReader(0, 2, false, false)
//
// 	res := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), test.NewReader(1, 1, false, true))
//
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return data.NewIntSet(1), res.AsSet(), nil
// 	})
//
// 	cp, rs, err := combinator.Optional(p).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, data.NewIntSet(1), cp)
// 	assert.Equal(t, parser.NewResultSet(res), rs)
// 	assert.Nil(t, err)
// }
//
// func TestOptionalShouldReturnEmptyResultIfParserFailed(t *testing.T) {
// 	r := test.NewReader(0, 2, false, false)
//
// 	p := parser.Func(func(leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return data.NewIntSet(1), nil, nil
// 	})
//
// 	cp, rs, err := combinator.Optional(p).Parse(parser.NewHistory(), parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, data.NewIntSet(1), cp)
// 	assert.Equal(t, parser.NewResultSet(parser.NewResult(nil, r)), rs)
// 	assert.Nil(t, err)
// }
