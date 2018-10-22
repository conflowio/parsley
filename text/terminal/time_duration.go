// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal

import (
	"errors"
	"time"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

// TimeDuration parses a duration string. A duration string is a possibly signed sequence of decimal numbers,
// each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func TimeDuration() parser.Func {
	notFoundErr := errors.New("was expecting a valid time duration")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*text.Reader)
		if readerPos, result := tr.ReadRegexp(pos, "[-+]?(?:[0-9]+(?:\\.[0-9]+)?(?:ns|us|µs|μs|ms|s|m|h))+"); result != nil {
			duration, err := time.ParseDuration(string(result))
			if err != nil {
				return nil, data.EmptyIntSet, parsley.NewError(pos, err)
			}
			return ast.NewTerminalNode("TIME_DURATION", duration, pos, readerPos), data.EmptyIntSet, nil
		}
		return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
	})
}
