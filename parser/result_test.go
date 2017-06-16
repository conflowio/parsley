package parser_test

import (
	"fmt"
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/test"
	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	r := test.NewReader(0, 0, false, false)
	node := ast.NewTerminalNode("x", r.Cursor(), "x")
	res := parser.NewResult(node, r)
	assert.Equal(t, node, res.Node())
	assert.Equal(t, r, res.Reader())
	assert.NotEmpty(t, res.String())
}

func testResult(pos int) parser.Result {
	r := test.NewReader(pos, 1, false, false)
	node := ast.NewTerminalNode("x", test.NewPosition(0), "x")
	return parser.NewResult(node, r)
}

func TestResultSetAppend(t *testing.T) {
	type TC struct {
		name     string
		actual   parser.ResultSet
		res      []parser.Result
		expected parser.ResultSet
	}

	testCases := []TC{
		TC{"Append to nil list", parser.ResultSet(nil), []parser.Result{testResult(1)}, parser.NewResultSet(testResult(1))},
		TC{"Append nil to list", parser.NewResultSet(testResult(1)), nil, parser.NewResultSet(testResult(1))},
		TC{"Append empty array to list", parser.NewResultSet(testResult(1)), []parser.Result{}, parser.NewResultSet(testResult(1))},
		TC{"Append to empty list", parser.NewResultSet(), []parser.Result{testResult(1)}, parser.NewResultSet(testResult(1))},
		TC{"Append to end of list", parser.NewResultSet(testResult(1)), []parser.Result{testResult(2)}, parser.NewResultSet(testResult(1), testResult(2))},
		TC{"Append should keep list sorted", parser.NewResultSet(testResult(2)), []parser.Result{testResult(1)}, parser.NewResultSet(testResult(1), testResult(2))},
		TC{"Append should ignore result with same cursor", parser.NewResultSet(testResult(2)), []parser.Result{testResult(2)}, parser.NewResultSet(testResult(2))},
		TC{"Append multiple", parser.NewResultSet(testResult(2)), []parser.Result{testResult(1), testResult(3)}, parser.NewResultSet(testResult(1), testResult(2), testResult(3))},
	}
	for _, tc := range testCases {
		tc.actual.Append(tc.res...)
		assert.Equal(t, tc.expected, tc.actual, fmt.Sprintf("Failed: %s", tc.name))
	}
}
