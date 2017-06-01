package terminal_test

import (
	"testing"

	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegerShouldMatch(t *testing.T) {
	type TC struct {
		input    string
		expected int
		cursor   int
	}
	testCases := []TC{
		TC{"0", 0, 1},
		TC{"+0", 0, 2},
		TC{"-0", 0, 2},
		TC{"1", 1, 1},
		TC{"+1", 1, 2},
		TC{"-1", -1, 2},
		TC{"1234567890", 1234567890, 10},
		TC{"+1234567890", 1234567890, 11},
		TC{"-1234567890", -1234567890, 11},
		TC{"123abc", 123, 3},
		TC{"00", 0, 1}, // 00 is not a valid integer, only "0" should be consumed
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res := terminal.Integer().Parse(parser.EmptyLeftRecCtx(), r)
		require.NotNil(t, res)
		actual, _ := res[0].Node().Value()
		assert.Equal(t, tc.expected, actual)
		assert.Equal(t, tc.cursor, res[0].Reader().Cursor().Pos())
	}
}

func TestIntegerShouldNotMatch(t *testing.T) {
	type TC struct {
		input string
	}
	testCases := []TC{
		TC{""},
		TC{" "},
		TC{"a"},
		TC{"-"},
		TC{"+"},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res := terminal.Integer().Parse(parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res)
	}
}
