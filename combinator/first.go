package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// First applies the given parsers until the first produces a result
func First(parsers ...parser.Parser) parser.Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		cp := parser.NoCurtailingParsers()
		var rs parser.ResultSet
		for _, p := range parsers {
			parser.Stat.RegisterCall()
			cp2, rs2 := p.Parse(leftRecCtx, r.Clone())
			cp = cp.Union(cp2)
			rs.Append(rs2...)
			if len(rs2) > 0 {
				break
			}
		}
		return cp, rs
	})
}
