package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// SuppressError removes the error from the parser result
func SuppressError(p parser.Parser) parser.Parser {
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		cp, rs, _ := p.Parse(leftRecCtx, r)
		return cp, rs, nil
	})
}
