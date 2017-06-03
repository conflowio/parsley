package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Recursive is a recursive and-type combinator
type Recursive struct {
	nodeBuilder       ast.NodeBuilder
	parserLookUp      func(i int) parser.Parser
	curtailingParsers data.IntSet
	resultSet         parser.ResultSet
	nodes             []ast.Node
	infinite          bool
}

// NewRecursive creates a new recursive combinator
func NewRecursive(nodeBuilder ast.NodeBuilder, infinite bool, parserLookUp func(i int) parser.Parser) *Recursive {
	return &Recursive{
		nodeBuilder:       nodeBuilder,
		infinite:          infinite,
		parserLookUp:      parserLookUp,
		curtailingParsers: parser.NoCurtailingParsers(),
		resultSet:         parser.ResultSet{},
		nodes:             []ast.Node{},
	}
}

// Parse runs the recursive parser
func (rp *Recursive) Parse(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
	rp.parse(0, 0, leftRecCtx, r, true)
	return rp.curtailingParsers, rp.resultSet
}

func (rp *Recursive) parse(depth int, nodesDepth int, leftRecCtx data.IntMap, r reader.Reader, mergeCurtailingParsers bool) bool {
	var cp data.IntSet
	var rs parser.ResultSet
	nextParser := rp.parserLookUp(depth)
	if nextParser != nil {
		parser.Stat.RegisterCall()
		cp, rs = nextParser.Parse(leftRecCtx, r.Clone())
	}

	if mergeCurtailingParsers {
		rp.curtailingParsers = rp.curtailingParsers.Union(cp)
	}

	if len(rs) > 0 {
		for i, result := range rs {
			nextNodeDepth := nodesDepth
			if result.Node() != nil {
				if len(rp.nodes) < nodesDepth+1 {
					rp.nodes = append(rp.nodes, result.Node())
				} else {
					rp.nodes[nodesDepth] = result.Node()
				}
				nextNodeDepth++
			}
			if i > 0 || result.Reader().Cursor().Pos() > r.Cursor().Pos() {
				leftRecCtx = parser.EmptyLeftRecCtx()
				mergeCurtailingParsers = false
			}
			if rp.parse(depth+1, nextNodeDepth, leftRecCtx, result.Reader().Clone(), mergeCurtailingParsers) {
				return true
			}
		}
	}
	if len(rs) == 0 {
		if (rp.infinite && depth > 0) || nextParser == nil {
			if nodesDepth > 0 {
				nodesCopy := make([]ast.Node, nodesDepth)
				copy(nodesCopy[0:nodesDepth], rp.nodes[0:nodesDepth])
				rp.resultSet.Append(parser.NewResult(rp.nodeBuilder.BuildNode(nodesCopy), r))
				if rp.nodes[nodesDepth-1].Token() == ast.EOF {
					return true
				}
			} else { // It's an empty result
				rp.resultSet.Append(parser.NewResult(nil, r))
			}
		}
	}
	return false
}
