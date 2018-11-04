// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// Evaluate parses the given input and evaluates it. It expects a reader, the root parser and the evaluation context.
// If there are multiple possible parse trees only the first one is used for evaluation.
func Evaluate(ctx *Context, p Parser, evalCtx interface{}) (interface{}, error) {
	node, parseErr := Parse(ctx, p)
	if parseErr != nil {
		return nil, parseErr
	}

	switch n := node.(type) {
	case StaticCheckable:
		if err := n.StaticCheck(evalCtx); err != nil {
			return nil, ctx.FileSet().ErrorWithPosition(err)
		}
	}

	value, evalErr := node.Value(evalCtx)
	if evalErr != nil {
		return nil, ctx.FileSet().ErrorWithPosition(evalErr)
	}

	return value, nil
}
