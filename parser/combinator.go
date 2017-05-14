package parser

import (
	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/reader"
)

// And combines multiple parsers
func And(nodeBuilder ast.NodeBuilder, parsers ...Parser) Func {
	return Func(func(r *reader.Reader) (ast.Node, *reader.Reader) {
		var nodes = make([]ast.Node, 0, len(parsers))
		var node ast.Node
		newReader := r.Clone()
		for _, parser := range parsers {
			node, newReader = parser.Parse(newReader)
			if node == nil {
				return nil, r
			}
			nodes = append(nodes, node)
		}
		return nodeBuilder(nodes), newReader
	})
}

// Or chooses the first matching parser
func Or(parsers ...Parser) Func {
	return Func(func(r *reader.Reader) (ast.Node, *reader.Reader) {
		for _, parser := range parsers {
			node, newReader := parser.Parse(r.Clone())
			if node != nil {
				return node, newReader
			}
		}
		return nil, nil
	})
}
