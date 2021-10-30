// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import (
	"github.com/conflowio/parsley/data"
)

// Parser defines a parser interface
//counterfeiter:generate . Parser
type Parser interface {
	Parse(ctx *Context, leftRecCtx data.IntMap, pos Pos) (Node, data.IntSet, Error)
}
