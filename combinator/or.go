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
	return parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		parserResult := parser.NewParserResult(parser.NoCurtailingParsers())
		for _, p := range parsers {
			parser.Stat.RegisterCall()
			r := p.Parse(leftRecCtx, r.Clone())
			if r != nil {
				parserResult.Append(r.Results...)
				parserResult.CurtailingParsers = parserResult.CurtailingParsers.Union(r.CurtailingParsers)
			}
		}
		return parserResult
	})
}
