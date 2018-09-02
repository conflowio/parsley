package parsley

import (
	"github.com/opsidian/parsley/data"
)

// Parser defines a parser interface
//go:generate counterfeiter . Parser
type Parser interface {
	Parse(ctx *Context, leftRecCtx data.IntMap, pos Pos) (Node, data.IntSet)
	Name() string
}

// Result is a stored parser result
type Result struct {
	LeftRecCtx        data.IntMap
	CurtailingParsers data.IntSet
	Node              Node
}

// History records information about parser calls
//go:generate counterfeiter . History
type History interface {
	SaveResult(parserIndex int, pos Pos, result *Result)
	GetResult(parserIndex int, pos Pos, leftRecCtx data.IntMap) (*Result, bool)
	RegisterCall()
	CallCount() int
}
