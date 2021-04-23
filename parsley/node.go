// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley

// Node represents an AST node
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Node
type Node interface {
	Token() string
	Schema() interface{}
	Pos() Pos
	ReaderPos() Pos
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . LiteralNode
type LiteralNode interface {
	Node
	Value() interface{}
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . NonLiteralNode
type NonLiteralNode interface {
	Node
	Value(userCtx interface{}) (interface{}, Error)
}

// NonTerminalNode represents a nonterminal AST node
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . NonTerminalNode
type NonTerminalNode interface {
	NonLiteralNode
	Children() []Node
}
