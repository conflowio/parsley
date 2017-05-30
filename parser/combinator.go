package parser

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/reader"
)

// Memoize handles result cache and curtailing left recursion
func Memoize(name string, c *Context, p Parser) Func {
	parserIndex := c.GetParserIndex(name)
	return Func(func(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
		result, found := c.GetResults(parserIndex, r.Cursor().Pos(), leftRecCtx)
		if found {
			return result
		}

		if leftRecCtx.Get(parserIndex) > r.CharsRemaining()+1 {
			return NewParserResult(data.NewIntSet().Insert(parserIndex))
		}

		result = p.Parse(leftRecCtx.Inc(parserIndex), r)
		if result != nil {
			leftRecCtx = leftRecCtx.Filter(result.CurtailingParsers)
		} else {
			leftRecCtx = data.NewIntMap()
		}

		c.RegisterResults(parserIndex, r.Cursor().Pos(), result, leftRecCtx)

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
			if r != nil {
				parserResult.Append(r.Results...)
				parserResult.CurtailingParsers = parserResult.CurtailingParsers.Union(r.CurtailingParsers)
			}
		}
		return parserResult
	}))
}

func parserListLookUp(parsers []Parser) func(i int) Parser {
	return func(i int) Parser {
		if i < len(parsers) {
			return parsers[i]
		}
		return nil
	}
}

// And combines multiple parsers
func And(name string, c *Context, nodeBuilder ast.NodeBuilder, parsers ...Parser) Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	return Memoize(name, c, Func(func(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
		return NewRecursiveParser(name, c, nodeBuilder, false, parserListLookUp(parsers)).Parse(leftRecCtx, r)
	}))
}

// Many matches the same expression one or more times
func Many(name string, c *Context, nodeBuilder ast.NodeBuilder, p Parser) Func {
	return Memoize(name, c, Func(func(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
		return NewRecursiveParser(name, c, nodeBuilder, true, func(i int) Parser { return p }).Parse(leftRecCtx, r)
	}))
}

// ManySep matches the given value parser one or more times separated by the separator parser
func ManySep(name string, token string, c *Context, valueP Parser, sepP Parser, interpreter ast.Interpreter) Func {
	sepValue := And(name+"_SV", c, ast.SingleNodeBuilder(1), sepP, valueP)
	sepValueMany := And(name+"_SV*", c, ast.AllNodesBuilder(token, interpreter), sepValue)
	return And(name, c, ast.AllNodesBuilder(token, interpreter), valueP, sepValueMany)
}

// RecursiveParser is a recursive and-type parser
type RecursiveParser struct {
	name         string
	c            *Context
	nodeBuilder  ast.NodeBuilder
	parserLookUp func(i int) Parser
	result       *ParserResult
	nodes        []ast.Node
	infinite     bool
}

// NewRecursiveParser creates a new recursive parser
func NewRecursiveParser(name string, c *Context, nodeBuilder ast.NodeBuilder, infinite bool, parserLookUp func(i int) Parser) RecursiveParser {
	return RecursiveParser{
		name:         name,
		c:            c,
		nodeBuilder:  nodeBuilder,
		infinite:     infinite,
		parserLookUp: parserLookUp,
		result:       NewParserResult(data.NewIntSet()),
		nodes:        []ast.Node{},
	}
}

// Parse runs the recursive parser
func (rp RecursiveParser) Parse(leftRecCtx data.IntMap, r *reader.Reader) *ParserResult {
	rp.runNextParser(0, leftRecCtx, r, true)
	return rp.result
}

func (rp RecursiveParser) runNextParser(depth int, leftRecCtx data.IntMap, r *reader.Reader, mergeCurtailingParsers bool) bool {
	var parserResult *ParserResult
	nextParser := rp.parserLookUp(depth)
	if nextParser != nil {
		rp.c.RegisterCall()
		parserResult = nextParser.Parse(leftRecCtx, r.Clone())
	}

	if parserResult != nil {
		if mergeCurtailingParsers {
			rp.result.CurtailingParsers = rp.result.CurtailingParsers.Union(parserResult.CurtailingParsers)
		}

		for i, result := range parserResult.Results {
			if len(rp.nodes) < depth+1 {
				rp.nodes = append(rp.nodes, result.Node())
			} else {
				rp.nodes[depth] = result.Node()
			}
			if i > 0 || result.Reader().Cursor().Pos() > r.Cursor().Pos() {
				leftRecCtx = data.NewIntMap()
				mergeCurtailingParsers = false
			}
			if rp.runNextParser(depth+1, leftRecCtx, result.Reader().Clone(), mergeCurtailingParsers) {
				return true
			}
		}
	}
	if parserResult == nil || parserResult.Results == nil {
		if (rp.infinite && depth > 0) || nextParser == nil {
			nodesCopy := make([]ast.Node, depth)
			copy(nodesCopy[0:depth], rp.nodes[0:depth])
			newRes := NewResult(rp.nodeBuilder.BuildNode(nodesCopy), r)
			rp.result.Append(newRes)
			if rp.nodes[depth-1].Token() == reader.EOF {
				return true
			}
		}
	}
	return false
}
