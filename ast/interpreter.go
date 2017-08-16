// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast

import "github.com/opsidian/parsley/reader"

// Interpreter defines an interface to evaluate AST nodes
type Interpreter interface {
	Eval(ctx interface{}, nodes []Node) (interface{}, reader.Error)
}

// InterpreterFunc defines a helper to implement the Interpreter interface with functions
type InterpreterFunc func(ctx interface{}, nodes []Node) (interface{}, reader.Error)

// Eval evaluates the given nodes and returns with a single result.
func (f InterpreterFunc) Eval(ctx interface{}, nodes []Node) (interface{}, reader.Error) {
	return f(ctx, nodes)
}
