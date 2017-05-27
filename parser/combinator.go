package parser

import (
	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/data"
	"github.com/opsidian/parsec/reader"
)

// Memoize handles result cache and curtailing left recursion
func Memoize(name string, c *Context, parser Func) Func {
	parserIndex := c.GetParserIndex(name)
	return Func(func(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
		result, found := c.GetResults(parserIndex, r.Cursor().Pos(), leftRecCtx)
		if found {
			return result
		}

		if leftRecCtx.Get(parserIndex) > r.CharsRemaining()+1 {
			return NewParserResult(data.NewIntSet().Insert(parserIndex))
		}

		result = parser(leftRecCtx.Inc(parserIndex), r)

		c.RegisterResults(parserIndex, r.Cursor().Pos(), result, leftRecCtx.Filter(result.CurtailingParsers))

		return result
	})
}

// Or chooses the first matching parser
func Or(name string, c *Context, parsers ...Parser) Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	return Memoize(name, c, Func(func(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
		parserResult := NewParserResult(data.NewIntSet())
		for _, parser := range parsers {
			c.RegisterCall()
			r := parser.Parse(leftRecCtx, r.Clone())
			parserResult.Append(r.Results...)
			parserResult.CurtailingParsers = parserResult.CurtailingParsers.Union(r.CurtailingParsers)
		}
		return parserResult
	}))
}

// And combines multiple parsers
func And(name string, c *Context, nodeBuilder ast.NodeBuilder, parsers ...Parser) Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	return Memoize(name, c, Func(func(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
		nodes := make([]ast.Node, len(parsers))
		result := NewParserResult(data.NewIntSet())
		andRec(c, leftRecCtx, nodeBuilder, result, 0, nodes, r, true, parsers...)
		return result
	}))
}

func andRec(c *Context, leftRecCtx data.IntMap, nodeBuilder ast.NodeBuilder, parserResult *ParserResult, depth int, nodes []ast.Node, r *reader.Reader, mergeCurtailingParsers bool, parsers ...Parser) bool {
	c.RegisterCall()
	nextParserResult := parsers[0].Parse(leftRecCtx, r.Clone())

	if nextParserResult != nil {
		if mergeCurtailingParsers {
			parserResult.CurtailingParsers = parserResult.CurtailingParsers.Union(nextParserResult.CurtailingParsers)
		}

		for i, result := range nextParserResult.Results {
			nodes[depth] = result.Node()
			if len(parsers) > 1 {
				var newLeftRecCtx data.IntMap
				var newMergeCurtailingParsers bool
				if i == 0 && result.Reader().Cursor().Pos() == r.Cursor().Pos() {
					newLeftRecCtx = leftRecCtx
					newMergeCurtailingParsers = true
				} else {
					newLeftRecCtx = data.NewIntMap()
					newMergeCurtailingParsers = false
				}
				if andRec(c, newLeftRecCtx, nodeBuilder, parserResult, depth+1, nodes, result.Reader().Clone(), newMergeCurtailingParsers, parsers[1:]...) {
					return true
				}
			} else {
				nodesCopy := make([]ast.Node, len(nodes))
				copy(nodesCopy, nodes)
				parserResult.Append(NewResult(nodeBuilder(nodesCopy), result.Reader()))
				if result.Node().Token() == reader.EOF {
					return true
				}
			}
		}
	}

	return false
}
