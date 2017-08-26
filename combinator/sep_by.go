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
	return newSepBy(token, valueP, sepP, interpreter, 0, false).CreateParser()
}

// SepByOrValue applies the given value parser zero or more times separated by the separator parser
// If there is only one value then the value node will be returned and the interpreter won't be used
// It simply uses the Seq, SeqTry, Many and Memoize combinators.
func SepByOrValue(token string, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Func {
	return newSepBy(token, valueP, sepP, interpreter, 0, true).CreateParser()
}

// SepBy1 applies the given value parser one or more times separated by the separator parser
// It simply uses the Seq, SeqTry, Many and Memoize combinators.
func SepBy1(token string, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Parser {
	return newSepBy(token, valueP, sepP, interpreter, 1, false).CreateParser()
}

// SepByOrValue1 applies the given value parser one or more times separated by the separator parser
// If there is only one value then the value node will be returned and the interpreter won't be used
// It simply uses the Seq, SeqTry, Many and Memoize combinators.
func SepByOrValue1(token string, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter) parser.Parser {
	return newSepBy(token, valueP, sepP, interpreter, 1, true).CreateParser()
}

type sepBy struct {
	token       string
	valueP      parser.Parser
	sepP        parser.Parser
	interpreter ast.Interpreter
	min         int
	returnValue bool
}

func newSepBy(token string, valueP parser.Parser, sepP parser.Parser, interpreter ast.Interpreter, min int, returnValue bool) sepBy {
	return sepBy{
		token:       token,
		valueP:      valueP,
		sepP:        sepP,
		interpreter: interpreter,
		min:         min,
		returnValue: returnValue,
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
		if s.returnValue && len(node1.Children()) == 0 {
			return nodes[0]
		}
		children = append(children, node1.Children()...)
	}
	return ast.NewNonTerminalNode(s.token, children, s.interpreter)
}
