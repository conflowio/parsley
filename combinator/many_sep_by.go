package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// ManySepBy applies the given value parser zero or more times separated by the separator parser
func ManySepBy(token string, h *parser.History, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Parser {
	return manySepBy{token, h, valueP, sepP, interpreter, 0}
}

// ManySepBy1 applies the given value parser one or more times separated by the separator parser
func ManySepBy1(token string, h *parser.History, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Parser {
	return manySepBy{token, h, valueP, sepP, interpreter, 1}
}

type manySepBy struct {
	token       string
	h           *parser.History
	valueP      parser.Parser
	sepP        parser.Parser
	interpreter ast.Interpreter
	min         int
}

func (m manySepBy) Parse(leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet) {
	sepValue := Memoize(m.token+"_MSB", m.h, And(builder.All("SEP_VALUE", nil), m.sepP, m.valueP))
	sepValueMany := Memoize(m.token+"_MSB*", m.h, Many(builder.Flatten(m.token, nil), sepValue))
	return Try(m, m.min, m.valueP, sepValueMany).Parse(leftRecCtx, r)
}

func (m manySepBy) BuildNode(nodes []ast.Node) ast.Node {
	if len(nodes) == 0 {
		return ast.NewNonTerminalNode(m.token, nil, m.interpreter)
	}

	children := []ast.Node{nodes[0]}
	if len(nodes) > 1 && nodes[1] != nil {
		node1 := nodes[1].(ast.NonTerminalNode)
		children = append(children, node1.Children()...)
	}
	return ast.NewNonTerminalNode(m.token, children, m.interpreter)
}
