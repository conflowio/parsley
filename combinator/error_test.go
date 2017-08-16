// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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

func TestSuppressErrorShouldReturnParserResultWithoutError(t *testing.T) {
	r := test.NewReader(0, 2, false, false)

	res := parser.NewResult(ast.NewTerminalNode("CHAR", test.NewPosition(1), 'a'), test.NewReader(1, 1, false, true))

	p := parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
		return data.NewIntSet(1), res.AsSet(), reader.NewError(test.NewPosition(1), "ERR1")
	})

	cp, rs, err := combinator.SuppressError(p).Parse(parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, data.NewIntSet(1), cp)
	assert.Equal(t, parser.NewResultSet(res), rs)
	assert.Nil(t, err)
}
