// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package parsley contains helper functions for parsing and evaluating text input data.
//
// Parsley is a general parser combinator library which can be used to parse context-free, left-recursive languages. It handles indirect as well as direct left-recursion in polynomial time and defines a memoization helper for speeding up parsing time. The language grammar can be easily translated to a set of rules using parsers and combinators.
//
// You can read a general introduction about parser combinators here: https://en.wikipedia.org/wiki/Parser_combinator
//
// For more information about handling left-recursion please check out **Parser Combinators for Ambiguous Left-Recursive Grammars (2007)** by Frost R.A., Hafiz R., and Callaghan P.
//
// Currently the library supports only text processing, but the interfaces are written with binary parsing in mind.
package parsley

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/text"
)

// ParseText parses the given text input and returns with the root node of the AST. It expects a byte array as input, the root parser and whether to ignore whitespaces between terminals. Unless you need full control over whitespaces in your parser set the ignoreWhitespaces to true.
// If there are multiple possible parse trees only the first one is returned.
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
// If there are multiple possible parse trees only the first one is used for evaluation.
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
