package parsley

import (
	"github.com/opsidian/parsley/data"
)

// Parser defines a parser interface
//go:generate counterfeiter . Parser
type Parser interface {
	Parse(h History, leftRecCtx data.IntMap, r Reader, pos int) (data.IntSet, []Node, Error)
}

// ParserFunc defines a helper to implement the Parser interface with functions
type ParserFunc func(h History, leftRecCtx data.IntMap, r Reader, pos int) (data.IntSet, []Node, Error)

// Parse parses the next token and returns with an AST node and the updated reader
func (p ParserFunc) Parse(h History, leftRecCtx data.IntMap, r Reader, pos int) (data.IntSet, []Node, Error) {
	return p(h, leftRecCtx, r, pos)
}

// Result is a stored parser result
type Result struct {
	LeftRecCtx        data.IntMap
	CurtailingParsers data.IntSet
	Nodes             []Node
	Err               Error
}

// History records information about parser calls
//go:generate counterfeiter . History
type History interface {
	SaveResult(parserIndex int, pos int, result *Result)
	GetResult(parserIndex int, pos int, leftRecCtx data.IntMap) (*Result, bool)
	RegisterCall()
	CallCount() int
}
