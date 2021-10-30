// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// Node represents an AST node
//counterfeiter:generate . Node
type Node interface {
	Token() string
	Schema() interface{}
	Pos() Pos
	ReaderPos() Pos
}

//counterfeiter:generate . LiteralNode
type LiteralNode interface {
	Node
	Value() interface{}
}

//counterfeiter:generate . NonLiteralNode
type NonLiteralNode interface {
	Node
	Value(userCtx interface{}) (interface{}, Error)
}

// NonTerminalNode represents a nonterminal AST node
//counterfeiter:generate . NonTerminalNode
type NonTerminalNode interface {
	NonLiteralNode
	Children() []Node
}
