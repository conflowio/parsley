package parser

import (
	"go/token"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/reader"
)

// And combines multiple parsers
func And(name string, nodeBuilder ast.NodeBuilder, parsers ...Parser) Func {
	return Func(func(h *History, r *reader.Reader) (results Results) {
		if parsers == nil {
			return nil
		}

		results, found := h.GetResults(name, r.Cursor().Pos())
		if found {
			return
		}

		if h.GetCalls(name, r.Cursor().Pos()) > r.CharsRemaining() {
			return nil
		}

		h.Push(name, r.Cursor().Pos())
		defer h.Pop(name, r.Cursor().Pos())

		nodes := make([]ast.Node, len(parsers))
		andRec(h, nodeBuilder, &results, 0, nodes, r, parsers...)
		h.RegisterResults(name, r.Cursor().Pos(), results)
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
			if result.Node().Token() != token.EOF {
				results.Add(nodeBuilder(nodesCopy), result.Reader())
			} else {
				*results = NewResults(Result{nodeBuilder(nodesCopy), result.Reader()})
				return true
			}
		}
	}
	return false
}

// Or chooses the first matching parser
func Or(name string, parsers ...Parser) Func {
	return Func(func(h *History, r *reader.Reader) (results Results) {
		if parsers == nil {
			return
		}

		results, found := h.GetResults(name, r.Cursor().Pos())
		if found {
			return
		}

		for _, parser := range parsers {
			for _, result := range parser.Parse(h, r.Clone()) {
				if result.Node().Token() != token.EOF {
					results = append(results, result)
				} else {
					return NewResults(result)
				}
			}
		}
		h.RegisterResults(name, r.Cursor().Pos(), results)
		return results
	})
}
