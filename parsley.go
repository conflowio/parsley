package parsley

import (
	"errors"

	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Parse parses the given input and returns with the root node of the AST
func Parse(input []byte, ignoreWhitespaces bool, s parser.Parser) (ast.Node, error) {
	parser.Stat.Reset()
	r := reader.New(input, ignoreWhitespaces)
	parser.Stat.RegisterCall()
	parserResult := s.Parse(parser.EmptyLeftRecCtx(), r)
	if len(parserResult.Results) == 0 {
		return nil, errors.New("failed to parse the input")
	}
	return parserResult.Results[0].Node(), nil
}

// Evaluate parses the given input and evaluates the AST
func Evaluate(input []byte, ignoreWhitespaces bool, s parser.Parser) (interface{}, error) {
	node, err := Parse(input, ignoreWhitespaces, s)
	if err != nil {
		return nil, err
	}
	return node.Value()
}
