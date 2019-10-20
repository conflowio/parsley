// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package text

import (
	"errors"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// LeftTrim skips the whitespaces before it tries to match the given parser
func LeftTrim(p parsley.Parser, wsMode WsMode) parser.Func {
	notFoundErr := errors.New("was expecting a new line")

	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		if wsMode == WsNone {
			return p.Parse(ctx, leftRecCtx, pos)
		}

		pos, ok := ctx.Reader().(*Reader).SkipWhitespaces(pos, wsMode)
		if !ok {
			return nil, data.EmptyIntSet, parsley.NewError(pos, notFoundErr)
		}
		return p.Parse(ctx, leftRecCtx, pos)
	})
}

// RightTrim reads and skips the whitespaces after any parser matches and updates the reader position
func RightTrim(p parsley.Parser, wsMode WsMode) parser.Func {
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		if wsMode == WsNone {
			return p.Parse(ctx, leftRecCtx, pos)
		}

		tr := ctx.Reader().(*Reader)
		res, cp, err := p.Parse(ctx, leftRecCtx, pos)
		if err != nil {
			errPos, _ := tr.SkipWhitespaces(err.Pos(), wsMode)
			if errPos > err.Pos() {
				err = parsley.NewError(errPos, err.Cause())
			}
			return res, cp, err
		}

		if res != nil {
			var ok bool
			res = ast.SetReaderPos(res, func(pos parsley.Pos) parsley.Pos {
				pos, ok = tr.SkipWhitespaces(pos, wsMode)
				return pos
			})
			if !ok {
				return nil, data.EmptyIntSet, parsley.NewError(res.ReaderPos(), errors.New("was expecting a new line"))
			}
		}

		return res, cp, nil
	})
}

// Trim removes all whitespaces before and after the result token
func Trim(p parsley.Parser) parser.Func {
	return RightTrim(LeftTrim(p, WsSpacesNl), WsSpacesNl)
}
