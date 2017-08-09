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

func TestIntegerShouldReturnStartingPosition(t *testing.T) {
	r := text.NewReader([]byte("1"), true)
	_, res, err := terminal.Integer().Parse(parser.EmptyLeftRecCtx(), r)
	require.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), res[0].Node().Pos())
}

func TestIntegerShouldMatch(t *testing.T) {
	type TC struct {
		input    string
		expected int
		cursor   int
	}
	testCases := []TC{
		{"0", 0, 1},
		{"+0", 0, 2},
		{"-0", 0, 2},
		{"1", 1, 1},
		{"+1", 1, 2},
		{"-1", -1, 2},
		{"1234567890", 1234567890, 10},
		{"+1234567890", 1234567890, 11},
		{"-1234567890", -1234567890, 11},
		{"123abc", 123, 3},
		{"00", 0, 2}, // this is a valid octal number
		{"01234567", 01234567, 8},
		{"+012", 012, 4},
		{"-012", -012, 4},
		{"08", 0, 1}, // as 08 is not a valid octal number only 0 should be parsed
		{"0x0123456789abcdef", 0x0123456789abcdef, 18},
		{"0X0123456789abcdef", 0x0123456789abcdef, 18},
		{"+0x12", 0x12, 5},
		{"-0x12", -0x12, 5},
		{"0xg", 0, 1}, // as 0xg is not a valid hexadecimal number only 0 should be parsed
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.Integer().Parse(parser.EmptyLeftRecCtx(), r)
		require.NotNil(t, res, fmt.Sprintf("Failed to parse %s", tc.input))
		actual, _ := res[0].Node().Value(nil)
		assert.Equal(t, tc.expected, actual)
		assert.Equal(t, tc.cursor, res[0].Reader().Cursor().Pos())
		assert.Nil(t, err)
	}
}

func TestIntegerShouldNotMatch(t *testing.T) {
	type TC struct {
		input string
	}
	testCases := []TC{
		{""},
		{"a"},
		{"-"},
		{"+"},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.Integer().Parse(parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res)
		require.NotNil(t, err)
		assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
	}
}
