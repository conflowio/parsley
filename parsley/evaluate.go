// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import (
	"errors"
)

var ErrNoValue = errors.New("node does not have a value")

// EvaluateNode evaluates the value of a node
// If the node doesn't have a value, it returns with an ErrNoValue error
func EvaluateNode(ctx interface{}, node Node) (interface{}, Error) {
	switch n := node.(type) {
	case LiteralNode:
		return n.Value(), nil
	case NonLiteralNode:
		return n.Value(ctx)
	default:
		return nil, NewError(node.Pos(), ErrNoValue)
	}
}

// Evaluate parses the given input and evaluates it. It expects a reader, the root parser and the evaluation context.
// If there are multiple possible parse trees only the first one is used for evaluation.
func Evaluate(ctx *Context, p Parser) (interface{}, error) {
	node, parseErr := Parse(ctx, p)
	if parseErr != nil {
		return nil, parseErr
	}

	value, evalErr := EvaluateNode(ctx.UserContext(), node)
	if evalErr != nil {
		return nil, ctx.FileSet().ErrorWithPosition(evalErr)
	}

	return value, nil
}
