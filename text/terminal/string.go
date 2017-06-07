package terminal

import (
	"strconv"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/token"
)

// String matches a string literal enclosed in double quotes
func String() parser.Func {
	return parser.Func(func(ctx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
		tr := r.(*text.Reader)
		matches, _, ok := tr.ReadMatch("\"|`")
		if !ok {
			return parser.NoCurtailingParsers(), nil
		}
		quote := matches[0]

		var value string
		var pos reader.Position
		if quote == "`" {
			var matches []string
			matches, pos, _ = tr.ReadMatch("[^`]*")
			value = matches[0]
		} else {
			value, pos, _ = tr.Readf(unquoteString)
		}

		endQuote, _, err := tr.ReadRune()
		if err != nil || string(endQuote) != quote {
			return parser.NoCurtailingParsers(), nil
		}
		return parser.NoCurtailingParsers(), parser.NewResult(ast.NewTerminalNode(token.STRING, pos, value), tr).AsSet()
	})
}

func unquoteString(b []byte) (string, int, bool) {
	str := string(b)
	var tail, res string
	var err error
	var ch rune
	for {
		if str == "" {
			break
		}
		ch, _, tail, err = strconv.UnquoteChar(str, '"')
		if err != nil {
			break
		}
		res += string(ch)
		str = tail
	}
	return res, len(b) - len(str), true
}
