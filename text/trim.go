// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package text

import (
	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
)

// LeftTrim skips the whitespaces before it tries to match the given parser
func LeftTrim(p parsley.Parser, wsMode WsMode) parser.Func {
	return func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*Reader)

		originalPos := pos

		pos, wsErr := tr.SkipWhitespaces(pos, wsMode)

		res, cp, err := p.Parse(ctx, leftRecCtx, pos)

		if ctxErr := ctx.Error(); ctxErr != nil {
			if ctxErr.Pos() == pos && parsley.IsNotFoundError(ctxErr) {
				ctx.SetError(parsley.NewError(originalPos, ctxErr.Cause()))
			}
		}

		if err != nil {
			if wsErr != nil {
				if err.Pos() > pos {
					return nil, data.EmptyIntSet, wsErr
				}

				if parsley.IsNotFoundError(err) {
					return res, cp, parsley.NewError(originalPos, err.Cause())
				}
			}

			return res, cp, err
		}

		if wsErr != nil {
			return nil, data.EmptyIntSet, wsErr
		}

		return res, cp, nil
	}
}

// RightTrim reads and skips the whitespaces after any parser matches and updates the reader position
func RightTrim(p parsley.Parser, wsMode WsMode) parser.Func {
	return func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
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
			var wsErr parsley.Error
			res = ast.SetReaderPos(res, func(pos parsley.Pos) parsley.Pos {
				pos, wsErr = tr.SkipWhitespaces(pos, wsMode)
				return pos
			})
			if wsErr != nil {
				return nil, data.EmptyIntSet, wsErr
			}
		}

		return res, cp, nil
	}
}

// Trim removes all whitespaces before and after the result token
func Trim(p parsley.Parser) parser.Func {
	return RightTrim(LeftTrim(p, WsSpacesNl), WsSpacesNl)
}
