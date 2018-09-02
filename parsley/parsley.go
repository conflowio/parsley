// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import (
	"fmt"

	"github.com/opsidian/parsley/data"
)

// Parse parses the given input and returns with the root node of the AST. It expects a reader and the root parser.
// If there are multiple possible parse trees only the first one is returned.
func Parse(parseCtx *Context, p Parser) (Node, Error) {
	parseCtx.History().RegisterCall()
	node, _ := p.Parse(parseCtx, data.EmptyIntMap, parseCtx.Reader().Pos(0))
	if node == nil {
		if parseCtx.Error() != nil {
			return nil, WrapError(parseCtx.Error(), "failed to parse the input: {{err}}")
		}
		return nil, NewError(parseCtx.Reader().Pos(0), fmt.Errorf("failed to parse the input: was expecting %s", p.Name()))
	}
	return node, nil
}

// Evaluate parses the given input and evaluates it. It expects a reader, the root parser and the evaluation context.
// If there are multiple possible parse trees only the first one is used for evaluation.
func Evaluate(parseCtx *Context, p Parser, evalCtx interface{}) (interface{}, Error) {
	node, parseErr := Parse(parseCtx, p)
	if parseErr != nil {
		return nil, parseErr
	}
	value, evalErr := node.Value(evalCtx)
	return value, evalErr
}
