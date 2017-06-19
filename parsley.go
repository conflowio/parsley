package parsley

import (
	"fmt"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
)

// ParseText parses the given text input and returns with the root node of the AST
func ParseText(input []byte, ignoreWhitespaces bool, s parser.Parser) (ast.Node, error) {
	parser.Stat.Reset()
	r := text.NewReader(input, ignoreWhitespaces)
	parser.Stat.RegisterCall()
	_, resultSet, err := s.Parse(parser.EmptyLeftRecCtx(), r)
	if len(resultSet) == 0 {
		return nil, fmt.Errorf("Failed to parse the input, %s", err.String())
	}
	if resultSet[0].Node() == nil {
		return nil, nil
	}
	return resultSet[0].Node(), nil
}

// EvaluateText parses the given text input and evaluates the AST
func EvaluateText(input []byte, ignoreWhitespaces bool, s parser.Parser) (interface{}, error) {
	node, err := ParseText(input, ignoreWhitespaces, s)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, nil
	}
	return node.Value()
}
