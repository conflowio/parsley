// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// Evaluate parses the given input and evaluates it. It expects a reader, the root parser and the evaluation context.
// If there are multiple possible parse trees only the first one is used for evaluation.
func Evaluate(ctx *Context, p Parser) (interface{}, error) {
	node, parseErr := Parse(ctx, p)
	if parseErr != nil {
		return nil, parseErr
	}

	value, evalErr := node.Value(ctx.UserContext())
	if evalErr != nil {
		return nil, ctx.FileSet().ErrorWithPosition(evalErr)
	}

	return value, nil
}
