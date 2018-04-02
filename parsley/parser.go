package parsley

import (
	"github.com/opsidian/parsley/data"
)

// Parser defines a parser interface
//go:generate counterfeiter . Parser
type Parser interface {
	Parse(h History, leftRecCtx data.IntMap, r Reader, pos int) (data.IntSet, []Node, Error)
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
	SaveResult(parserIndex int, pos int, result Result)
	GetResult(parserIndex int, pos int, leftRecCtx data.IntMap) (Result, bool)
	RegisterCall()
	CallCount() int
}
