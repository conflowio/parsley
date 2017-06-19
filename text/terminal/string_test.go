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

func TestStringShouldMatch(t *testing.T) {
	type TC struct {
		input    string
		expected string
		cursor   int
	}
	testCases := []TC{
		TC{`""`, "", 2},
		TC{`"a"`, "a", 3},
		TC{`" a "`, " a ", 5}, // Should keep whitespaces in the string
		TC{`"a" "b"`, "a", 3},
		TC{`"abcd"`, "abcd", 6},
		TC{` "abcd"`, "abcd", 7}, // Skip whitespaces
		TC{`"'"`, "'", 3},
		TC{`"\a\b\f\n\r\t\v"`, "\a\b\f\n\r\t\v", 16},
		TC{`"\x67"`, "\x67", 6},
		TC{`"\uAB12"`, "\uAB12", 8},
		TC{`"\U0001F355"`, "\U0001F355", 12},
		TC{"``", ``, 2},
		TC{"`a`", `a`, 3},
		TC{"` a `", ` a `, 5}, // Should keep whitespaces in the string
		TC{"`a` `b`", `a`, 3},
		TC{"`abcd`", `abcd`, 6},
		TC{" `abcd`", `abcd`, 7}, // Skip whitespaces
		TC{"`'`", `'`, 3},
		TC{"`" + `\a\b\f\n\r\t\v` + "`", `\a\b\f\n\r\t\v`, 16},
		TC{"`" + `\x67` + "`", `\x67`, 6},
		TC{"`" + `\uAB12` + "`", `\uAB12`, 8},
		TC{"`" + `\U0001F355` + "`", `\U0001F355`, 12},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.String().Parse(parser.EmptyLeftRecCtx(), r)
		require.NotNil(t, res, fmt.Sprintf("Failed to parse: %s", tc.input))
		actual, _ := res[0].Node().Value()
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
		TC{``},
		TC{`'`},
		TC{"''"},
		TC{"'a'"},
		TC{"5"},
		TC{"a"},
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
		TC{`"`, text.NewPosition(1, 1, 2)},
		TC{"`", text.NewPosition(1, 1, 2)},
		TC{`"a`, text.NewPosition(2, 1, 3)},
		TC{"`a", text.NewPosition(2, 1, 3)},
	}
	for _, tc := range testCases {
		r := text.NewReader([]byte(tc.input), true)
		_, res, err := terminal.String().Parse(parser.EmptyLeftRecCtx(), r)
		require.Nil(t, res, fmt.Sprintf("Should fail to parse: %s", tc.input))
		require.NotNil(t, err)
		assert.Equal(t, tc.errPos, err.Pos())
	}
}
