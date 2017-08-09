// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal_test

import (
	"fmt"
	"testing"

	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCharShouldReturnQuotePosition(t *testing.T) {
	r := text.NewReader([]byte("'a'"), true)
	_, res, err := terminal.Char().Parse(parser.EmptyLeftRecCtx(), r)
	require.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), res[0].Node().Pos())
}

func TestCharShouldMatch(t *testing.T) {
	type TC struct {
		input    string
		expected rune
		cursor   int
	}
	testCases := []TC{
		{`' '`, ' ', 3},
		{`'a'`, 'a', 3},
		{`'🍕'`, '🍕', 6},
		{`'a' xx`, 'a', 3},
		{`'\a'`, '\a', 4},
		{`'\b'`, '\b', 4},
		{`'\f'`, '\f', 4},
		{`'\n'`, '\n', 4},
		{`'\v'`, '\v', 4},
		{`'\r'`, '\r', 4},
		{`'\t'`, '\t', 4},
		{`'\xa0'`, '\xa0', 6},
		{`'\xA0'`, '\xa0', 6},
		{`'\u12ab'`, '\u12ab', 8},
		{`'\u12AB'`, '\u12ab', 8},
		{`'\U0001f355'`, '🍕', 12},
		{`'\U0001F355'`, '🍕', 12},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.Char().Parse(parser.EmptyLeftRecCtx(), r)
		require.NotNil(t, res, fmt.Sprintf("Failed to parse: %s", tc.input))
		actual, _ := res[0].Node().Value(nil)
		assert.Equal(t, tc.expected, actual)
		assert.Equal(t, tc.cursor, res[0].Reader().Cursor().Pos())
		assert.Nil(t, err)
	}
}

func TestCharShouldNotMatch(t *testing.T) {
	type TC struct {
		input string
	}
	testCases := []TC{
		{``},
		{`''`},
		{`'aa'`},
		{`"a"`},
		{`\x`},
		{`\x`},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.Char().Parse(parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res)
		require.NotNil(t, err)
		assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
	}
}
