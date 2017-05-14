package parser

import (
	"github.com/opsidian/parsec/ast"
	"github.com/opsidian/parsec/reader"
)

// Parser defines a parser interface
type Parser interface {
	Parse(*reader.Reader) (ast.Node, *reader.Reader)
}

// Func defines a helper to implement the Parser interface with functions
type Func func(*reader.Reader) (ast.Node, *reader.Reader)

// Parse parses the next token and returns with an AST node and the updated reader
func (f Func) Parse(r *reader.Reader) (ast.Node, *reader.Reader) {
	return f(r)
}
