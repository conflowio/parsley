package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// recursive is a recursive and-type combinator
type recursive struct {
	nodeBuilder       ast.NodeBuilder
	parserLookUp      func(i int) parser.Parser
	min               int
	max               int
	curtailingParsers data.IntSet
	resultSet         parser.ResultSet
	nodes             []ast.Node
}

// newRecursive creates a new recursive combinator
func newRecursive(nodeBuilder ast.NodeBuilder, parserLookUp func(i int) parser.Parser, min int, max int) *recursive {
	if nodeBuilder == nil {
		panic("Node builder can not be nil!")
	}
	return &recursive{
		nodeBuilder:       nodeBuilder,
		parserLookUp:      parserLookUp,
		min:               min,
		max:               max,
		curtailingParsers: parser.NoCurtailingParsers(),
		resultSet:         parser.ResultSet{},
		nodes:             []ast.Node{},
	}
}

// Parse runs the recursive parser
func (rp *recursive) Parse(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
	rp.parse(0, leftRecCtx, r, true)
	return rp.curtailingParsers, rp.resultSet
}

func (rp *recursive) parse(depth int, leftRecCtx data.IntMap, r reader.Reader, mergeCurtailingParsers bool) bool {
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
			if len(rp.nodes) < depth+1 {
				rp.nodes = append(rp.nodes, result.Node())
			} else {
				rp.nodes[depth] = result.Node()
			}
			if i > 0 || result.Reader().Cursor().Pos() > r.Cursor().Pos() {
				leftRecCtx = parser.EmptyLeftRecCtx()
				mergeCurtailingParsers = false
			}
			if rp.parse(depth+1, leftRecCtx, result.Reader().Clone(), mergeCurtailingParsers) {
				return true
			}
		}
	}
	if len(rs) == 0 {
		if depth >= rp.min && (rp.max == -1 || depth <= rp.max) {
			if depth > 0 {
				nodesCopy := make([]ast.Node, depth)
				copy(nodesCopy[0:depth], rp.nodes[0:depth])
				rp.resultSet.Append(parser.NewResult(rp.nodeBuilder.BuildNode(nodesCopy), r))
				if rp.nodes[depth-1] != nil && rp.nodes[depth-1].Token() == ast.EOF {
					return true
				}
			} else { // It's an empty result
				rp.resultSet.Append(parser.NewResult(rp.nodeBuilder.BuildNode(nil), r))
			}
		}
	}
	return false
}
