package parser

import (
	"fmt"

	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/reader"
)

// Parser defines a parser interface
type Parser interface {
	Parse(*History, *reader.Reader) Results
}

// Func defines a helper to implement the Parser interface with functions
type Func func(*History, *reader.Reader) Results

// Parse parses the next token and returns with an AST node and the updated reader
func (f Func) Parse(h *History, r *reader.Reader) Results {
	return f(h, r)
}

// Result represents one result of a parser
type Result struct {
	node   ast.Node
	reader *reader.Reader
}

// Node returns with the node
func (r Result) Node() ast.Node {
	return r.node
}

// Reader returns with the reader
func (r Result) Reader() *reader.Reader {
	return r.reader
}

func (r Result) String() string {
	return fmt.Sprintf("%s, cur: %s", r.node, r.reader.Cursor())
}

// Results is a result array
type Results []Result

// NewResults creates a new results instance
func NewResults(results ...Result) Results {
	return Results(results)
}

// Add adds a new result
func (r *Results) Add(node ast.Node, reader *reader.Reader) {
	*r = append(*r, Result{node, reader})
}
