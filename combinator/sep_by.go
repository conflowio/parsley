package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// SepBy applies the given value parser zero or more times separated by the separator parser
// It simply uses the Seq, SeqTry, Many and Memoize combinators.
func SepBy(token string, h *parser.History, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Parser {
	return sepBy{token, h, valueP, sepP, interpreter, 0}
}

// SepBy1 applies the given value parser one or more times separated by the separator parser
// It simply uses the Seq, SeqTry, Many and Memoize combinators.
func SepBy1(token string, h *parser.History, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Parser {
	return sepBy{token, h, valueP, sepP, interpreter, 1}
}

type sepBy struct {
	token       string
	h           *parser.History
	valueP      parser.Parser
	sepP        parser.Parser
	interpreter ast.Interpreter
	min         int
}

func (s sepBy) Parse(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, parser.Error) {
	sepValue := Memoize(s.token+"_MSB", s.h, Seq(builder.All("SEP_VALUE", nil), s.sepP, s.valueP))
	sepValueMany := Memoize(s.token+"_MSB*", s.h, Many(builder.Flatten(s.token, nil), sepValue))
	return SeqTry(s, s.min, s.valueP, sepValueMany).Parse(leftRecCtx, r)
}

func (s sepBy) BuildNode(nodes []ast.Node) ast.Node {
	if len(nodes) == 0 {
		return ast.NewNonTerminalNode(s.token, nil, s.interpreter)
	}

	children := []ast.Node{nodes[0]}
	if len(nodes) > 1 && nodes[1] != nil {
		node1 := nodes[1].(ast.NonTerminalNode)
		children = append(children, node1.Children()...)
	}
	return ast.NewNonTerminalNode(s.token, children, s.interpreter)
}
