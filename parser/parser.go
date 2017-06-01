package parser

import (
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/reader"
)

// Parser defines a parser interface
type Parser interface {
	Parse(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, ResultSet)
}

// Func defines a helper to implement the Parser interface with functions
type Func func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, ResultSet)

// Parse parses the next token and returns with an AST node and the updated reader
func (f Func) Parse(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, ResultSet) {
	return f(leftRecCtx, r)
}

// EmptyLeftRecCtx creates an empty left recursion context
func EmptyLeftRecCtx() data.IntMap {
	return data.EmptyIntMap()
}

// NoCurtailingParsers returns with an empty int set
func NoCurtailingParsers() data.IntSet {
	return data.EmptyIntSet()
}
