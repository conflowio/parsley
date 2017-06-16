package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Or tries all the given parsers and merges the results
func Or(parsers ...parser.Parser) parser.Func {
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
		}
		return cp, rs
	})
}
