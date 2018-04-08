package text

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// LeftTrim skips the whitespaces before it tries to match the given parser
func LeftTrim(p parsley.Parser, wsMode WsMode) parser.Func {
	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, parsley.Node, parsley.Error) {
		return p.Parse(h, leftRecCtx, r, r.(*Reader).SkipWhitespaces(pos, wsMode))
	})
}

// RightTrim reads and skips the whitespaces after any parser matches and updates the reader position
func RightTrim(p parsley.Parser, wsMode WsMode) parser.Func {
	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, parsley.Node, parsley.Error) {
		tr := r.(*Reader)
		cp, res, err := p.Parse(h, leftRecCtx, r, pos)
		if res != nil {
			res.SetReaderPos(func(pos int) int { return tr.SkipWhitespaces(pos, wsMode) })
		}
		return cp, res, err
	})
}
