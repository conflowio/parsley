package terminal

import (
	"fmt"
	"strconv"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/token"
)

// Float matches a float literal
func Float(bitSize int) parser.Func {
	return parser.Func(func(ctx data.IntMap, r reader.Reader) (cp data.IntSet, rs parser.ResultSet) {
		tr := r.(*text.Reader)
		if matches, pos := tr.ReadMatch("^[-+]?[0-9]*\\.?[0-9]+(?:[eE][-+]?[0-9]+)?"); matches != nil {
			val, err := strconv.ParseFloat(matches[0], bitSize)
			if err != nil {
				panic(fmt.Sprintf("Invalid float value encountered: %s", matches[0]))
			}
			if bitSize == 32 {
				rs = parser.NewResult(ast.NewTerminalNode(token.FLOAT, pos, float32(val)), r).AsSet()
			} else {
				rs = parser.NewResult(ast.NewTerminalNode(token.FLOAT, pos, val), r).AsSet()
			}
		}
		return
	})
}
