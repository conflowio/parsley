package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Choice tries to apply the given parsers until one of them succeeds
func Choice(parsers ...parser.Parser) parser.Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	return parser.Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
		cp := parser.NoCurtailingParsers()
		var rs parser.ResultSet
		var err parser.Error
		for _, p := range parsers {
			parser.Stat.RegisterCall()
			cp2, rs2, err2 := p.Parse(leftRecCtx, r.Clone())
			cp = cp.Union(cp2)
			rs.Append(rs2...)
			if len(rs2) > 0 {
				break
			}
			if err2 != nil && (err == nil || err2.Pos().Pos() >= err.Pos().Pos()) {
				err = err2
			}
		}
		if len(rs) > 0 {
			return cp, rs, nil
		} else {
			return cp, nil, err
		}
	})
}
