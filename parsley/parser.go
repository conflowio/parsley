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
