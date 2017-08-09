// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ast defines structs and interfaces for building and evaluating an abstract syntax tree.
//
// Abstract Syntax Tree
//
// ASTs can have two types of nodes. A non-terminal node (branch node) can contain a token and multiple child nodes.
// A terminal node (leaf node) can contain a token and a value.
//
// Builders
//
// Builder are used to build a node from multiple nodes.
//
// Interpreters
//
// Interpreters get a list of nodes and return with a single value.
package ast
