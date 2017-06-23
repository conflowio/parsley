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
		TC{`' '`, ' ', 3},
		TC{`'a'`, 'a', 3},
		TC{`'🍕'`, '🍕', 6},
		TC{`'a' xx`, 'a', 3},
		TC{`'\a'`, '\a', 4},
		TC{`'\b'`, '\b', 4},
		TC{`'\f'`, '\f', 4},
		TC{`'\n'`, '\n', 4},
		TC{`'\v'`, '\v', 4},
		TC{`'\r'`, '\r', 4},
		TC{`'\t'`, '\t', 4},
		TC{`'\xa0'`, '\xa0', 6},
		TC{`'\xA0'`, '\xa0', 6},
		TC{`'\u12ab'`, '\u12ab', 8},
		TC{`'\u12AB'`, '\u12ab', 8},
		TC{`'\U0001f355'`, '🍕', 12},
		TC{`'\U0001F355'`, '🍕', 12},
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
		TC{``},
		TC{`''`},
		TC{`'aa'`},
		TC{`"a"`},
		TC{`\x`},
		TC{`\x`},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.Char().Parse(parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res)
		require.NotNil(t, err)
		assert.Equal(t, text.NewPosition(0, 1, 1), err.Pos())
	}
}
