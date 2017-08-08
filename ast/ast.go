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
