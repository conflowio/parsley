// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/builder"
	"github.com/opsidian/parsley/parser"
)

// SepBy applies the given value parser zero or more times separated by the separator parser
// It simply uses the Seq, SeqTry, Many and Memoize combinators.
func SepBy(token string, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Func {
	return newSepBy(token, valueP, sepP, interpreter, 0).CreateParser()
}

// SepBy1 applies the given value parser one or more times separated by the separator parser
// It simply uses the Seq, SeqTry, Many and Memoize combinators.
func SepBy1(token string, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Parser {
	return newSepBy(token, valueP, sepP, interpreter, 1).CreateParser()
}

type sepBy struct {
	token       string
	valueP      parser.Parser
	sepP        parser.Parser
	interpreter ast.Interpreter
	min         int
}

func newSepBy(token string, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter, min int) sepBy {
	return sepBy{
		token:       token,
		valueP:      valueP,
		sepP:        sepP,
		interpreter: interpreter,
		min:         min,
	}
}

func (s sepBy) CreateParser() parser.Func {
	sepValue := Memoize(Seq(builder.All("SEP_VALUE", nil), s.sepP, s.valueP))
	sepValueMany := Memoize(Many(builder.Flatten(s.token, nil), sepValue))
	return SeqTry(s, s.min, s.valueP, sepValueMany)
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
