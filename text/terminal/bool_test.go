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

func TestBoolShouldReturnStartingPosition(t *testing.T) {
	r := text.NewReader([]byte("true"), "", true)
	_, res, err := terminal.Bool().Parse(nil, parser.EmptyLeftRecCtx(), r)
	require.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), res[0].Node().Pos())
}

func TestBoolShouldMatch(t *testing.T) {
	type TC struct {
		input    string
		expected bool
		cursor   int
	}
	testCases := []TC{
		{"true", true, 4},
		{"false", false, 5},
		{"true xxx", true, 4},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), "", true)
		_, res, err := terminal.Bool().Parse(nil, parser.EmptyLeftRecCtx(), r)
		require.NotNil(t, res, fmt.Sprintf("Failed to parse: %s", tc.input))
		actual, _ := res[0].Node().Value(nil)
		assert.Equal(t, tc.expected, actual)
		assert.Equal(t, tc.cursor, res[0].Reader().Cursor().Pos())
		assert.Nil(t, err)
	}
}

func TestBoolShouldNotMatch(t *testing.T) {
	type TC struct {
		input string
	}
	testCases := []TC{
		{"tru"},
		{"fals"},
		{"TRUE"},
		{"FALSE"},
		{"truex"},
		{"falsex"},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), "", true)
		_, res, err := terminal.Bool().Parse(nil, parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res)
		require.NotNil(t, err)
		assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
	}
}
