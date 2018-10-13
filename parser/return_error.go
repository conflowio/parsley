package parser

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

// ReturnError will override the returned error by the parser if its position is the same as the reader's position
func ReturnError(p parsley.Parser, customErr error) Func {
	return Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
		res, cp, err := p.Parse(ctx, leftRecCtx, pos)

		if err != nil {
			if err.Pos() == pos {
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
