package parser

import (
	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/reader"
)

// And combines multiple parsers
func And(name string, nodeBuilder ast.NodeBuilder, parsers ...Parser) Func {
	return Func(func(h *History, r *reader.Reader) Results {
		if parsers == nil {
			return nil
		}

		if pos, ok := h.GetLastPosition(name); ok && pos == r.Position() {
			return nil
		}
		h.Push(name, r.Position())
		defer h.Pop(name)

		results := NewResults()
		nodes := make([]ast.Node, len(parsers))
		andRec(h, nodeBuilder, &results, 0, nodes, r, parsers...)
		return results
	})
}

func andRec(h *History, nodeBuilder ast.NodeBuilder, results *Results, depth int, nodes []ast.Node, r *reader.Reader, parsers ...Parser) bool {
	for _, result := range parsers[0].Parse(h, r.Clone()) {
		nodes[depth] = result.Node()
		if len(parsers) > 1 {
			if andRec(h, nodeBuilder, results, depth+1, nodes, result.Reader().Clone(), parsers[1:]...) {
				return true
			}
		} else {
			nodesCopy := make([]ast.Node, len(nodes))
			copy(nodesCopy, nodes)
			results.Add(nodeBuilder(nodesCopy), result.Reader())
			// TODO: only check for eof
			if result.Reader().ReadEOF() {
				return true
			}
		}
	}
	return false
}

// Or chooses the first matching parser
func Or(name string, parsers ...Parser) Func {
	return Func(func(h *History, r *reader.Reader) Results {
		if parsers == nil {
			return nil
		}

		var results []Result
		for _, parser := range parsers {
			for _, result := range parser.Parse(h, r.Clone()) {
				results = append(results, result)
				// TODO: this should be just a check, not a read
				if result.Reader().ReadEOF() {
					break
				}
			}
		}
		return results
	})
}
