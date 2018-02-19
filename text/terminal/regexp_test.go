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

func TestRegexpShouldMatchPattern(t *testing.T) {
	r := text.NewReader([]byte("abc123"), "", true)
	_, res, err := terminal.Regexp("letters", "[a-z]+", false, 0, "TEST")(nil, parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("TEST", text.NewPosition(0, 1, 1), "abc")
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assert.Equal(t, text.NewPosition(3, 1, 4), r.Cursor())
	assert.Nil(t, err)
}

func TestRegexpShouldIgnoreWhitespaces(t *testing.T) {
	r := text.NewReader([]byte("   abc123"), "", true)
	_, res, err := terminal.Regexp("letters", "[a-z]+", false, 0, "TEST")(nil, parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("TEST", text.NewPosition(3, 1, 4), "abc")
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assert.Equal(t, text.NewPosition(6, 1, 7), r.Cursor())
	assert.Nil(t, err)
}

func TestRegexpShouldNotIgnoreWhitespaces(t *testing.T) {
	r := text.NewReader([]byte("   abc123a"), "", true)
	_, res, err := terminal.Regexp("letters", "[ a-z]+", true, 0, "TEST")(nil, parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("TEST", text.NewPosition(0, 1, 1), "   abc")
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assert.Equal(t, text.NewPosition(6, 1, 7), r.Cursor())
	assert.Nil(t, err)
}

func TestRegexpShouldUseMatchIfTokenEmpty(t *testing.T) {
	r := text.NewReader([]byte("abc123"), "", true)
	_, res, err := terminal.Regexp("letter", "[a-z]+", false, 0, "")(nil, parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("abc", text.NewPosition(0, 1, 1), "abc")
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assert.Equal(t, text.NewPosition(3, 1, 4), r.Cursor())
	assert.Nil(t, err)
}

func TestRegexpNotMatchingPattern(t *testing.T) {
	r := text.NewReader([]byte("abc123"), "", true)
	_, res, err := terminal.Regexp("numbers", "[0-9]+", false, 0, "TEST")(nil, parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
	require.NotNil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
}

func TestRegexpShouldPanicIfInvalidGroupIndex(t *testing.T) {
	r := text.NewReader([]byte("abc123"), "", true)
	assert.Panics(t, func() { terminal.Regexp("letters", "[a-z]+", false, 1, "TEST")(nil, parser.EmptyLeftRecCtx(), r) })
}
