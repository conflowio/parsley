package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/parser"
)

// ManySepBy matches the given value parser one or more times separated by the separator parser
func ManySepBy(name string, token string, h *parser.History, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Func {
	sepValue := Memoize(name+"_SV", h, And(builder.All("TMP", interpreter), sepP, valueP))
	sepValueMany := Memoize(name+"_SV*", h, Maybe(Many(builder.Flatten(token, interpreter), sepValue)))
	return And(builder.Flatten(token, interpreter), valueP, sepValueMany)
}
