// Package parsley is a parser combinator library.
//
// You can read a general introduction about parser combinators here: https://en.wikipedia.org/wiki/Parser_combinator
package parsley

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
)

// ParseText parses the given text input and returns with the root node of the AST. It expects a byte array as input, the root parser and whether to ignore whitespaces between terminals. Unless you need full control over whitespaces in your parser set the ignoreWhitespaces to true.
func ParseText(input []byte, ignoreWhitespaces bool, s parser.Parser) (ast.Node, parser.Error) {
	parser.Stat.Reset()
	r := text.NewReader(input, ignoreWhitespaces)
	parser.Stat.RegisterCall()
	_, resultSet, err := s.Parse(parser.EmptyLeftRecCtx(), r)
	if len(resultSet) == 0 {
		return nil, parser.NewError(err.Pos(), "Failed to parse the input, %s", err.Msg())
	}
	if resultSet[0].Node() == nil {
		return nil, nil
	}
	return resultSet[0].Node(), nil
}

// EvaluateText parses the given text input and evaluates it. It expects a byte array as input, the root parser, the evaluation context and whether to ignore whitespaces between terminals. Unless you need full control over whitespaces in your parser set the ignoreWhitespaces to true.
func EvaluateText(input []byte, ignoreWhitespaces bool, s parser.Parser, ctx interface{}) (interface{}, error) {
	node, err := ParseText(input, ignoreWhitespaces, s)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, nil
	}
	return node.Value(ctx)
}
