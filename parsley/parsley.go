// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import (
	"github.com/opsidian/parsley/data"
)

// Parse parses the given input and returns with the root node of the AST. It expects a reader and the root parser.
// If there are multiple possible parse trees only the first one is returned.
func Parse(ctx *Context, p Parser) (Node, error) {
	node, _, err := p.Parse(ctx, data.EmptyIntMap, ctx.Reader().Pos(0))
	if err != nil {
		return nil, ctx.FileSet().ErrorWithPosition(err)
	}

	return node, nil
}

// Evaluate parses the given input and evaluates it. It expects a reader, the root parser and the evaluation context.
// If there are multiple possible parse trees only the first one is used for evaluation.
func Evaluate(ctx *Context, p Parser, evalCtx interface{}) (interface{}, error) {
	node, parseErr := Parse(ctx, p)
	if parseErr != nil {
		return nil, parseErr
	}

	value, evalErr := node.Value(evalCtx)
	if evalErr != nil {
		return nil, ctx.FileSet().ErrorWithPosition(evalErr)
	}

	return value, nil
}
