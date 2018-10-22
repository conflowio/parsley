// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package parser contains the main structs for parsing
package parser

import (
	"fmt"

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
	return ReturnError(Func(f), fmt.Errorf("was expecting %s", name))
}
