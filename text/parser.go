package text

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
)

// LeftTrim skips the whitespaces before it tries to match the given parser
func LeftTrim(p parsley.Parser, wsMode WsMode) *parser.NamedFunc {
	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (parsley.Node, parsley.Error, data.IntSet) {
		return p.Parse(h, leftRecCtx, r, r.(*Reader).SkipWhitespaces(pos, wsMode))
	}).WithName(p.Name())
}

// RightTrim reads and skips the whitespaces after any parser matches and updates the reader position
func RightTrim(p parsley.Parser, wsMode WsMode) *parser.NamedFunc {
	return parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (parsley.Node, parsley.Error, data.IntSet) {
		tr := r.(*Reader)
		res, err, cp := p.Parse(h, leftRecCtx, r, pos)
		if res != nil {
			res = ast.SetReaderPos(res, func(pos parsley.Pos) parsley.Pos { return tr.SkipWhitespaces(pos, wsMode) })
		}
		return res, err, cp
	}).WithName(p.Name())
}
