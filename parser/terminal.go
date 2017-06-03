package parser

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/reader"
)

// Empty always matches and returns with an empty node result
func Empty() Func {
	return Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, ResultSet) {
		return NoCurtailingParsers(), NewResult(nil, r).AsSet()
	})
}

// End matches the end of the input
func End() Func {
	return Func(func(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, ResultSet) {
		if r.IsEOF() {
			return NoCurtailingParsers(), NewResult(ast.NewTerminalNode(ast.EOF, r.Cursor(), nil), r).AsSet()
		}
		return NoCurtailingParsers(), nil
	})
}
