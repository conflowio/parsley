// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package test

import (
	"testing"

	"github.com/opsidian/parsley/ast"
	"github.com/stretchr/testify/assert"
)

// AssertNodesEqual asserts that two AST nodes are equal
func AssertNodesEqual(t *testing.T, ctx interface{}, expected ast.Node, actual ast.Node) {
	switch n := expected.(type) {
	case ast.TerminalNode:
		assert.Equal(t, n, actual.(ast.TerminalNode))
	case ast.NonTerminalNode:
		n2 := actual.(ast.NonTerminalNode)
		assert.Equal(t, n.Token(), n2.Token())
		assert.Equal(t, n.Children(), n2.Children())
		val1, err1 := n.Value(ctx)
		val2, err2 := n2.Value(ctx)
		assert.Equal(t, val1, val2)
		assert.Equal(t, err1, err2)
	}
}
