// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator_test

// func TestSuppressErrorShouldReturnParserResultWithoutError(t *testing.T) {
// 	r := test.NewReader(0, 2, false, false)
//
// 	res := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), test.NewReader(1, 1, false, true))
//
// 	p := parsley.ParserFunc(func(leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
// 		return data.NewIntSet(1), res.AsSet(), reader.NewError(test.NewPosition(1), "ERR1")
// 	})
//
// 	h := parser.NewHistory()
// 	cp, rs, err := combinator.SuppressError(p).Parse(h, parser.EmptyLeftRecCtx(), r)
// 	assert.Equal(t, data.NewIntSet(1), cp)
// 	assert.Equal(t, parser.NewResultSet(res), rs)
// 	assert.Nil(t, err)
// 	assert.Equal(t, 0, h.CallCount())
// }
