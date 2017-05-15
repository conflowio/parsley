package parser

import (
	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/reader"
)

// And combines multiple parsers
func And(nodeBuilder ast.NodeBuilder, parsers ...Parser) Func {
	return Func(func(r *reader.Reader) Results {
		if parsers == nil {
			return nil
		}
		results := NewResults()
		nodes := make([]ast.Node, len(parsers))
		andRec(nodeBuilder, &results, 0, nodes, r, parsers...)
		return results
	})
}

func andRec(nodeBuilder ast.NodeBuilder, results *Results, depth int, nodes []ast.Node, r *reader.Reader, parsers ...Parser) {
	for _, result := range parsers[0].Parse(r.Clone()) {
		nodes[depth] = result.Node()
		if len(parsers) > 1 {
			andRec(nodeBuilder, results, depth+1, nodes, result.Reader().Clone(), parsers[1:]...)
		} else {
			nodesCopy := make([]ast.Node, len(nodes))
			copy(nodesCopy, nodes)
			results.Add(nodeBuilder(nodesCopy), result.Reader())
		}
	}
}

// Or chooses the first matching parser
func Or(parsers ...Parser) Func {
	return Func(func(r *reader.Reader) Results {
		if parsers == nil {
			return nil
		}
		var results []Result
		for _, parser := range parsers {
			for _, result := range parser.Parse(r.Clone()) {
				results = append(results, result)
			}
		}
		return results
	})
}
