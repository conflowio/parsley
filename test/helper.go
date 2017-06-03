package test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/stretchr/testify/assert"
)

// AssertNodesEqual asserts that two AST nodes are equal
func AssertNodesEqual(t *testing.T, expected ast.Node, actual ast.Node) {
	switch n := expected.(type) {
	case ast.TerminalNode:
		assert.Equal(t, n, actual.(ast.TerminalNode))
	case ast.NonTerminalNode:
		n2 := actual.(ast.NonTerminalNode)
		assert.Equal(t, n.Token(), n2.Token())
		assert.Equal(t, n.Children(), n2.Children())
		val1, err1 := n.Value()
		val2, err2 := n2.Value()
		assert.Equal(t, val1, val2)
		assert.Equal(t, err1, err2)
	}
}
