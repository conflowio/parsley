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
type Func func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (parsley.Node, parsley.Error, data.IntSet)

func (f Func) Parse(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (parsley.Node, parsley.Error, data.IntSet) {
	return f(h, leftRecCtx, r, pos)
}

func (f Func) Name() string {
	return ""
}

func (f Func) WithName(name string) *NamedFunc {
	return &NamedFunc{
		name: name,
		f:    f,
	}
}

type NamedFunc struct {
	name string
	f    Func
}

func (nf *NamedFunc) Parse(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (parsley.Node, parsley.Error, data.IntSet) {
	return nf.f(h, leftRecCtx, r, pos)
}

// Name returns with the parser name
func (nf *NamedFunc) Name() string {
	return nf.name
}
