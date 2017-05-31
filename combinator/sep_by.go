package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
)

// SepBy matches the given value parser one or more times separated by the separator parser
func SepBy(name string, token string, h *parser.History, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Func {
	sepValue := Memoize(name+"_SV", h, And(ast.SingleNodeBuilder(1), sepP, valueP))
	sepValueMany := Memoize(name+"_SV*", h, Many(ast.AllNodesBuilder(token, interpreter), sepValue))
	return And(ast.AllNodesBuilder(token, interpreter), valueP, sepValueMany)
}
