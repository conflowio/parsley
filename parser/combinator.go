package parser

import (
	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/reader"
)

// Memoize handles result cache and curtailing left recursion
func Memoize(name string, parser Func) Func {
	return Func(func(c *Context, r *reader.Reader) (results *Results) {
		c.Push(name, r.Cursor().Pos())
		defer c.Pop(name, r.Cursor().Pos())

		results, found := c.GetResults(name, r.Cursor().Pos())
		if found {
			return
		}

		if c.GetCalls(name, r.Cursor().Pos()) > r.CharsRemaining()+1 {
			return c.NewCurtailedResults(r.Cursor().Pos())
		}

		results = parser(c, r)
		c.RegisterResults(name, r.Cursor().Pos(), results)
		return
	})
}

// Or chooses the first matching parser
func Or(name string, parsers ...Parser) Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	return Memoize(name, Func(func(c *Context, r *reader.Reader) (results *Results) {
		results = NewResults(nil)
		for _, parser := range parsers {
			c.RegisterCall()
			if r := parser.Parse(c, r.Clone()); r != nil {
				results.Merge(r)
			}
		}
		return results
	}))
}

// And combines multiple parsers
func And(name string, nodeBuilder ast.NodeBuilder, parsers ...Parser) Func {
	if parsers == nil {
		panic("No parsers were given")
	}
	return Memoize(name, Func(func(c *Context, r *reader.Reader) (results *Results) {
		nodes := make([]ast.Node, len(parsers))
		results = NewResults(nil)
		andRec(c, nodeBuilder, results, 0, nodes, r, parsers...)
		return results
	}))
}

func andRec(c *Context, nodeBuilder ast.NodeBuilder, results *Results, depth int, nodes []ast.Node, r *reader.Reader, parsers ...Parser) bool {
	c.RegisterCall()
	results2 := parsers[0].Parse(c, r.Clone())
	if results2 != nil {
		if depth == 0 || len(results2.items) == 0 {
			results.MergeCurtailReasons(results2)
		}
		for _, result := range results2.items {
			nodes[depth] = result.Node()
			if len(parsers) > 1 {
				if andRec(c, nodeBuilder, results, depth+1, nodes, result.Reader().Clone(), parsers[1:]...) {
					return true
				}
			} else {
				nodesCopy := make([]ast.Node, len(nodes))
				copy(nodesCopy, nodes)
				results.Add(NewResult(nodeBuilder(nodesCopy), result.Reader()))
				if result.Node().Token() == reader.EOF {
					return true
				}
			}
		}
	}

	return false
}
