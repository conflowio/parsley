// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal_test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWordShouldMatchExactString(t *testing.T) {
	r := text.NewReader([]byte("foo bar"), "", true)
	_, res, err := terminal.Word("foo", "TEST", 5).Parse(nil, parser.EmptyLeftRecCtx(), r)
	assert.Equal(t, res, parser.NewResult(ast.NewTerminalNode("TEST", text.NewPosition(0, 1, 1), 5), r).AsSet())
	assert.Nil(t, err)
}

func TestWordShouldNotMatchWordPartially(t *testing.T) {
	r := text.NewReader([]byte("foobar"), "", true)
	_, res, err := terminal.Word("foo", "TEST", 5).Parse(nil, parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	require.NotNil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
}

func TestWordShouldMatchOnlyGivenString(t *testing.T) {
	r := text.NewReader([]byte("bar"), "", true)
	_, res, err := terminal.Word("foo", "TEST", 5).Parse(nil, parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	require.NotNil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
}
