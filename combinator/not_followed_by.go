package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// NotFollowedBy returns only with the parsers results which are not followed by the second parser results
func NotFollowedBy(p1 parser.Parser, p2 parser.Parser) parser.Func {
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		validRs := parser.NewResultSet()
		cp, rs := p1.Parse(leftRecCtx, r)
		for _, result := range rs {
			_, rs2 := p2.Parse(leftRecCtx, result.Reader().Clone())
			if len(rs2) == 0 {
				validRs.Append(result)
			}
		}
		return cp, validRs
	})
}
