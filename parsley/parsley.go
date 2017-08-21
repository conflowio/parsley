// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// Sentence is the root parser which processes the whole input
type Sentence struct {
	p parser.Parser
}

// NewSentence creates a new Sentence instance
func NewSentence(p parser.Parser) *Sentence {
	return &Sentence{
		p: combinator.Seq(builder.Select(0), p, parser.End()),
	}
}

// Parse will parse the whole input
func (s *Sentence) Parse(r reader.Reader) (ast.Node, *parser.History, reader.Error) {
	return Parse(r, s.p)
}

// Evaluate will parse the whole input and evaluate the AST
func (s *Sentence) Evaluate(r reader.Reader, ctx interface{}) (interface{}, *parser.History, reader.Error) {
	return Evaluate(r, s.p, ctx)
}

// Parse parses the given input and returns with the root node of the AST. It expects a reader and the root parser.
// If there are multiple possible parse trees only the first one is returned.
func Parse(r reader.Reader, s parser.Parser) (ast.Node, *parser.History, reader.Error) {
	h := parser.NewHistory()
	h.RegisterCall()
	_, resultSet, err := s.Parse(h, parser.EmptyLeftRecCtx(), r)
	if len(resultSet) == 0 {
		return nil, h, reader.WrapError(err.Pos(), err, "Failed to parse the input: %s", err.Msg())
	}
	if resultSet[0].Node() == nil {
		return nil, h, nil
	}
	return resultSet[0].Node(), h, nil
}

// Evaluate parses the given input and evaluates it. It expects a reader, the root parser and the evaluation context.
// If there are multiple possible parse trees only the first one is used for evaluation.
func Evaluate(r reader.Reader, s parser.Parser, ctx interface{}) (interface{}, *parser.History, reader.Error) {
	node, h, parseErr := Parse(r, s)
	if parseErr != nil {
		return nil, h, parseErr
	}
	if node == nil {
		return nil, h, nil
	}
	value, evalErr := node.Value(ctx)
	return value, h, evalErr
}
