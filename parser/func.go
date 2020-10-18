// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package parser contains the main structs for parsing
package parser

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

// Func defines a helper to implement the Parser interface with functions
type Func func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error)

// Parse parses the input using the function
func (f Func) Parse(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	return f(ctx, leftRecCtx, pos)
}

// Name returns with a new parser function which overrides the returned error
// if its position is the same as the reader's position
// The error will be: "was expecting <name>"
func (f Func) Name(name string) Func {
	return ReturnError(Func(f), parsley.NewNotFoundError(name))
}

// FuncWrapper is a parser which wraps a parser function as a struct
// It's useful when you have to use a parser recursively as that's not possible with functions
type FuncWrapper struct {
	F Func
}

// Parse parses the input using the function
func (f FuncWrapper) Parse(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	return f.F(ctx, leftRecCtx, pos)
}
