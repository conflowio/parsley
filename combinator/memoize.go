package combinator

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Memoize handles result cache and curtailing left recursion
func Memoize(name string, h *parser.History, p parser.Parser) parser.Func {
	parserIndex := h.GetParserIndex(name)
	return parser.Func(func(leftRecCtx data.IntMap, r *reader.Reader) *parser.ParserResult {
		result, found := h.GetResults(parserIndex, r.Cursor().Pos(), leftRecCtx)
		if found {
			return result
		}

		if leftRecCtx.Get(parserIndex) > r.CharsRemaining()+1 {
			return parser.NewParserResult(data.NewIntSet(parserIndex))
		}

		result = p.Parse(leftRecCtx.Inc(parserIndex), r)
		if result != nil {
			leftRecCtx = leftRecCtx.Filter(result.CurtailingParsers)
		} else {
			leftRecCtx = parser.EmptyLeftRecCtx()
		}

		h.RegisterResults(parserIndex, r.Cursor().Pos(), result, leftRecCtx)

		return result
	})
}
