// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import (
	"github.com/opsidian/parsley/data"
)

// Parser defines a parser interface
//go:generate counterfeiter . Parser
type Parser interface {
	Parse(ctx *Context, leftRecCtx data.IntMap, pos Pos) (Node, data.IntSet, Error)
}
