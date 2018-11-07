// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// StaticChecker defines an interface to run a static analysis on the given nonterminal node
//go:generate counterfeiter . StaticChecker
type StaticChecker interface {
	StaticCheck(ctx interface{}, node NonTerminalNode) (string, Error)
}

// StaticCheckable is an interface for nodes that can run a static analysis
type StaticCheckable interface {
	StaticCheck(ctx interface{}) Error
}

// StaticCheckableNode defines a node which also implements the StaticCheckable interface
//go:generate counterfeiter . StaticCheckableNode
type StaticCheckableNode interface {
	Node
	StaticCheck(ctx interface{}) Error
}

// StaticCheck parses the given input and runs static analysis on it. It expects a reader, the root parser and the evaluation context.
// If there are multiple possible parse trees only the first one is used for the analysis.
func StaticCheck(ctx *Context, node Node, evalCtx interface{}) error {
	switch n := node.(type) {
	case StaticCheckable:
		if err := n.StaticCheck(evalCtx); err != nil {
			return ctx.FileSet().ErrorWithPosition(err)
		}
	}

	return nil
}
