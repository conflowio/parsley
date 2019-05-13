// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// StaticChecker defines an interface to run a static analysis on the given nonterminal node
//go:generate counterfeiter . StaticChecker
type StaticChecker interface {
	StaticCheck(userCtx interface{}, node NonTerminalNode) (string, Error)
}

// StaticCheckable is an interface for nodes that can run a static analysis
type StaticCheckable interface {
	StaticCheck(userCtx interface{}) Error
}

// StaticCheckableNode defines a node which also implements the StaticCheckable interface
//go:generate counterfeiter . StaticCheckableNode
type StaticCheckableNode interface {
	Node
	StaticCheckable
}

// StaticCheck will run static checking on the given node
func StaticCheck(userCtx interface{}, node Node) Error {
	var staticCheckErr Error

	Walk(node, func(n Node) bool {
		if staticCheckableNode, ok := n.(StaticCheckable); ok {
			if err := staticCheckableNode.StaticCheck(userCtx); err != nil {
				staticCheckErr = err
				return true
			}
		}

		return false
	})

	return staticCheckErr
}
