// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parsley"
)

// ReturnError will override the returned error by the parser if its position is the same as the reader's position
func ReturnError(p parsley.Parser, customErr error) Func {
	return Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		res, cp, err := p.Parse(ctx, leftRecCtx, pos)

		if err != nil {
			if err.Pos() == pos && parsley.IsNotFoundError(err) {
				err = parsley.NewError(pos, customErr)
			}

			return nil, cp, err
		}

		if res == nil {
			err = parsley.NewError(pos, customErr)
		}

		return res, cp, err
	})
}
