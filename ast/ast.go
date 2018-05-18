// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ast defines structs and interfaces for building and evaluating an abstract syntax tree.
//
// Abstract Syntax Tree
//
// ASTs can have different types of nodes.
// * A terminal node (leaf node) can contain a token and a value.
// * A nil node (leaf node) only contains a position and evaluates to nil
// * A non-terminal node (branch node) can contain a token and multiple child nodes.
//
// Interpreters
//
// Interpreters get a list of nodes and return with a single value.
package ast
