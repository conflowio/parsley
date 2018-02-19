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

func TestRuneShouldMatchCharacter(t *testing.T) {
	r := text.NewReader([]byte("a"), "", true)
	_, res, err := terminal.Rune('a', "A")(nil, parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("A", text.NewPosition(0, 1, 1), 'a')
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assert.Equal(t, text.NewPosition(1, 1, 2), r.Cursor())
	assert.Nil(t, err)
}

func TestRuneShouldNotUseSpecialChars(t *testing.T) {
	r := text.NewReader([]byte("a"), "", true)
	_, res, err := terminal.Rune('.', ".")(nil, parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
	require.NotNil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
}

func TestRuneShouldMatchUnicodeCharacter(t *testing.T) {
	r := text.NewReader([]byte("🍕"), "", true)
	_, res, err := terminal.Rune('🍕', "PIZZA")(nil, parser.EmptyLeftRecCtx(), r)
	expectedNode := ast.NewTerminalNode("PIZZA", text.NewPosition(0, 1, 1), '🍕')
	assert.Equal(t, parser.NewResult(expectedNode, r).AsSet(), res)
	assert.Equal(t, text.NewPosition(4, 1, 2), r.Cursor())
	assert.Nil(t, err)
}

func TestRuneShouldNotMatchCharacter(t *testing.T) {
	r := text.NewReader([]byte("a"), "", true)
	_, res, err := terminal.Rune('b', "B")(nil, parser.EmptyLeftRecCtx(), r)
	assert.Nil(t, res)
	assert.Equal(t, text.NewPosition(0, 1, 1), r.Cursor())
	require.NotNil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
}
