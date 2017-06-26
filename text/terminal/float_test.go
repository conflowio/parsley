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
	TC{"0.1", 0.1, 3},
	TC{"+0.1", 0.1, 4},
	TC{"-0.1", -0.1, 4},
	TC{"1.1", 1.1, 3},
	TC{"+1.1", 1.1, 4},
	TC{"-1.1", -1.1, 4},
	TC{"1234567890.1234567890", 1234567890.1234567890, 21},
	TC{"+1234567890.1234567890", 1234567890.1234567890, 22},
	TC{"-1234567890.1234567890", -1234567890.1234567890, 22},
	TC{"123.4abc", 123.4, 5},
	TC{"00.00", 0, 5},
	TC{".0", 0, 2},
	TC{".1", 0.1, 2},
	TC{"+.1", 0.1, 3},
	TC{"-.1", -0.1, 3},
	TC{"1.2e5", 1.2e5, 5},
	TC{"+1.2e5", 1.2e5, 6},
	TC{"-1.2e5", -1.2e5, 6},
	TC{"1.2e", 1.2, 3}, // only 1.2 should be consumed
}

func TestFloatShouldReturnStartingPosition(t *testing.T) {
	r := text.NewReader([]byte("1.23"), true)
	_, res, err := terminal.Float().Parse(parser.EmptyLeftRecCtx(), r)
	require.NotNil(t, res)
	assert.Nil(t, err)
	assert.Equal(t, text.NewPosition(0, 1, 1), res[0].Node().Pos())
}

func TestFloatShouldMatch(t *testing.T) {
	for _, tc := range validTestCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.Float().Parse(parser.EmptyLeftRecCtx(), r)
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
		TC{""},
		TC{"a"},
		TC{"-"},
		TC{"+"},
		TC{"."},
		TC{"1"},
		TC{"-1"},
		TC{"0"},
		TC{"1.2e3456"},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.Float().Parse(parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res)
		require.NotNil(t, err)
		assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
	}
}
