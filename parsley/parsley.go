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
func Parse(h History, r Reader, s Parser) (Node, Error) {
	h.RegisterCall()
	node, err, _ := s.Parse(h, data.EmptyIntMap, r, Pos(1))
	if err != nil {
		return nil, WrapError(err, "Failed to parse the input: {{err}}")
	}
	if node == nil {
		return nil, NewError(r.Pos(0), "Failed to parse the input: was expecting %s", s.Name())
	}
	return node, nil
}

// Evaluate parses the given input and evaluates it. It expects a reader, the root parser and the evaluation context.
// If there are multiple possible parse trees only the first one is used for evaluation.
func Evaluate(h History, r Reader, s Parser, ctx interface{}) (interface{}, Error) {
	node, parseErr := Parse(h, r, s)
	if parseErr != nil {
		return nil, parseErr
	}
	if node == nil {
		return nil, nil
	}
	value, evalErr := node.Value(ctx)
	return value, evalErr
}
