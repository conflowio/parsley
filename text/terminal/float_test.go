// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal_test

import (
	"testing"

	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TC struct {
	input    string
	expected float64
	cursor   int
}

var validTestCases = []TC{
	{"0.1", 0.1, 3},
	{"+0.1", 0.1, 4},
	{"-0.1", -0.1, 4},
	{"1.1", 1.1, 3},
	{"+1.1", 1.1, 4},
	{"-1.1", -1.1, 4},
	{"1234567890.1234567890", 1234567890.1234567890, 21},
	{"+1234567890.1234567890", 1234567890.1234567890, 22},
	{"-1234567890.1234567890", -1234567890.1234567890, 22},
	{"123.4abc", 123.4, 5},
	{"00.00", 0, 5},
	{".0", 0, 2},
	{".1", 0.1, 2},
	{"+.1", 0.1, 3},
	{"-.1", -0.1, 3},
	{"1.2e5", 1.2e5, 5},
	{"+1.2e5", 1.2e5, 6},
	{"-1.2e5", -1.2e5, 6},
	{"1.2e", 1.2, 3}, // only 1.2 should be consumed
}

func TestFloatShouldReturnStartingPosition(t *testing.T) {
	r := text.NewReader([]byte("1.23"), "", true)
	_, res, err := terminal.Float().Parse(nil, parser.EmptyLeftRecCtx(), r)
	require.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), res[0].Node().Pos())
}

func TestFloatShouldMatch(t *testing.T) {
	for _, tc := range validTestCases {
		r := text.NewReader([]byte(tc.input), "", true)
		_, res, err := terminal.Float().Parse(nil, parser.EmptyLeftRecCtx(), r)
		require.NotNil(t, res)
		actual, _ := res[0].Node().Value(nil)
		assert.Equal(t, tc.expected, actual)
		assert.Equal(t, tc.cursor, res[0].Reader().Cursor().Pos())
		assert.Nil(t, err)
	}
}

func TestFloatShouldNotMatch(t *testing.T) {
	type TC struct {
		input string
	}
	testCases := []TC{
		{""},
		{"a"},
		{"-"},
		{"+"},
		{"."},
		{"1"},
		{"-1"},
		{"0"},
		{"1.2e3456"},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), "", true)
		_, res, err := terminal.Float().Parse(nil, parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res)
		require.NotNil(t, err)
		assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
	}
}
