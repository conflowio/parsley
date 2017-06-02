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

func TestBoolShouldMatch(t *testing.T) {
	type TC struct {
		input    string
		expected bool
		cursor   int
	}
	testCases := []TC{
		TC{"true", true, 4},
		TC{"false", false, 5},
		TC{"true xxx", true, 4},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res := terminal.Bool().Parse(parser.EmptyLeftRecCtx(), r)
		require.NotNil(t, res, fmt.Sprintf("Failed to parse: %s", tc.input))
		actual, _ := res[0].Node().Value()
		assert.Equal(t, tc.expected, actual)
		assert.Equal(t, tc.cursor, res[0].Reader().Cursor().Pos())
	}
}

func TestBoolShouldNotMatch(t *testing.T) {
	type TC struct {
		input string
	}
	testCases := []TC{
		TC{"tru"},
		TC{"fals"},
		TC{"TRUE"},
		TC{"FALSE"},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res := terminal.Bool().Parse(parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res)
	}
}
