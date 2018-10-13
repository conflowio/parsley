package text

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// LeftTrim skips the whitespaces before it tries to match the given parser
func LeftTrim(p parsley.Parser, wsMode WsMode) parser.Func {
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		pos = ctx.Reader().(*Reader).SkipWhitespaces(pos, wsMode)
		return p.Parse(ctx, leftRecCtx, pos)
	})
}

// RightTrim reads and skips the whitespaces after any parser matches and updates the reader position
func RightTrim(p parsley.Parser, wsMode WsMode) parser.Func {
	return parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		tr := ctx.Reader().(*Reader)
		res, cp, err := p.Parse(ctx, leftRecCtx, pos)
		if err != nil {
			errPos := tr.SkipWhitespaces(err.Pos(), wsMode)
			if errPos > err.Pos() {
				err = parsley.NewError(errPos, err.Cause())
			}
			return res, cp, err
		}

		if res != nil {
			res = ast.SetReaderPos(res, func(pos parsley.Pos) parsley.Pos { return tr.SkipWhitespaces(pos, wsMode) })
		}

		return res, cp, nil
	})
}

// Trim removes all whitespaces before and after the result token
func Trim(p parsley.Parser) parser.Func {
	return RightTrim(LeftTrim(p, WsSpacesNl), WsSpacesNl)
}
