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

func TestStringShouldReturnQuotePosition(t *testing.T) {
	r := text.NewReader([]byte(`"abc"`), true)
	_, res, err := terminal.String().Parse(parser.EmptyLeftRecCtx(), r)
	require.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), res[0].Node().Pos())
}

func TestStringShouldMatch(t *testing.T) {
	type TC struct {
		input    string
		expected string
		cursor   int
	}
	testCases := []TC{
		{`""`, "", 2},
		{`"a"`, "a", 3},
		{`" a "`, " a ", 5}, // Should keep whitespaces in the string
		{`"a" "b"`, "a", 3},
		{`"abcd"`, "abcd", 6},
		{` "abcd"`, "abcd", 7}, // Skip whitespaces
		{`"'"`, "'", 3},
		{`"\a\b\f\n\r\t\v"`, "\a\b\f\n\r\t\v", 16},
		{`"\x67"`, "\x67", 6},
		{`"\uAB12"`, "\uAB12", 8},
		{`"\U0001F355"`, "\U0001F355", 12},
		{"``", ``, 2},
		{"`a`", `a`, 3},
		{"` a `", ` a `, 5}, // Should keep whitespaces in the string
		{"`a` `b`", `a`, 3},
		{"`abcd`", `abcd`, 6},
		{" `abcd`", `abcd`, 7}, // Skip whitespaces
		{"`'`", `'`, 3},
		{"`" + `\a\b\f\n\r\t\v` + "`", `\a\b\f\n\r\t\v`, 16},
		{"`" + `\x67` + "`", `\x67`, 6},
		{"`" + `\uAB12` + "`", `\uAB12`, 8},
		{"`" + `\U0001F355` + "`", `\U0001F355`, 12},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.String().Parse(parser.EmptyLeftRecCtx(), r)
		require.NotNil(t, res, fmt.Sprintf("Failed to parse: %s", tc.input))
		actual, _ := res[0].Node().Value(nil)
		assert.Equal(t, tc.expected, actual)
		assert.Equal(t, tc.cursor, res[0].Reader().Cursor().Pos())
		assert.Nil(t, err)
	}
}

func TestStringShouldNotMatch(t *testing.T) {
	type TC struct {
		input string
	}
	testCases := []TC{
		{``},
		{`'`},
		{"''"},
		{"'a'"},
		{"5"},
		{"a"},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.String().Parse(parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res, fmt.Sprintf("Should fail to parse: %s", tc.input))
		require.NotNil(t, err)
		assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
	}
}

func TestUnclosedStringLiteral(t *testing.T) {
	type TC struct {
		input  string
		errPos text.Position
	}
	testCases := []TC{
		{`"`, text.NewPosition(1, 1, 2)},
		{"`", text.NewPosition(1, 1, 2)},
		{`"a`, text.NewPosition(2, 1, 3)},
		{"`a", text.NewPosition(2, 1, 3)},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.String().Parse(parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res, fmt.Sprintf("Should fail to parse: %s", tc.input))
		require.NotNil(t, err)
		assert.Equal(t, tc.errPos, err.Pos())
	}
}
