// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
)

// Recursive is a recursive and-type combinator
type Recursive struct {
	token        string
	parserLookUp func(int) parsley.Parser
	lenCheck     func(int) bool
	interpreter  parsley.Interpreter
	notFoundErr  error
}

// NewRecursive creates a new recursive instance
func NewRecursive(token string, parserLookUp func(int) parsley.Parser, lenCheck func(int) bool) *Recursive {
	return &Recursive{
		token:        token,
		parserLookUp: parserLookUp,
		lenCheck:     lenCheck,
	}
}

// Bind binds the given interpreter
func (rp *Recursive) Bind(interpreter parsley.Interpreter) *Recursive {
	rp.interpreter = interpreter
	return rp
}

// Parse parses the given input
func (rp *Recursive) Parse(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	p := &recursive{
		token:             rp.token,
		parserLookUp:      rp.parserLookUp,
		lenCheck:          rp.lenCheck,
		interpreter:       rp.interpreter,
		curtailingParsers: data.EmptyIntSet,
		nodes:             []parsley.Node{},
	}
	return p.Parse(ctx, leftRecCtx, pos)
}

// recursive is a recursive and-type combinator
type recursive struct {
	token             string
	parserLookUp      func(i int) parsley.Parser
	lenCheck          func(i int) bool
	interpreter       parsley.Interpreter
	curtailingParsers data.IntSet
	result            parsley.Node
	err               parsley.Error
	nodes             []parsley.Node
}

// Parse runs the recursive parser
func (rp *recursive) Parse(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
	rp.parse(0, ctx, leftRecCtx, pos, true)
	if rp.result == nil {
		return nil, rp.curtailingParsers, rp.err
	}

	return rp.result, rp.curtailingParsers, nil
}

func (rp *recursive) parse(depth int, ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos, mergeCurtailingParsers bool) bool {
	var cp data.IntSet
	var res parsley.Node
	var err parsley.Error
	nextParser := rp.parserLookUp(depth)
	if nextParser != nil {
		ctx.RegisterCall()
		res, cp, err = nextParser.Parse(ctx, leftRecCtx, pos)
		if err != nil && (rp.err == nil || err.Pos() > rp.err.Pos()) {
			rp.err = err
		}
	}

	if mergeCurtailingParsers {
		rp.curtailingParsers = rp.curtailingParsers.Union(cp)
	}

	if res != nil {
		switch rest := res.(type) {
		case ast.NodeList:
			for i, node := range rest {
				if rp.parseNext(i, node, depth, ctx, leftRecCtx, pos, mergeCurtailingParsers) {
					return true
				}
			}
		default:
			if rp.parseNext(0, rest, depth, ctx, leftRecCtx, pos, mergeCurtailingParsers) {
				return true
			}
		}
	}

	if res == nil {
		if rp.lenCheck(depth) {
			if depth > 0 {
				nodesCopy := make([]parsley.Node, depth)
				copy(nodesCopy[0:depth], rp.nodes[0:depth])
				rp.result = ast.AppendNode(rp.result, ast.NewNonTerminalNode(rp.token, nodesCopy, rp.interpreter))
				if rp.nodes[depth-1] != nil && rp.nodes[depth-1].Token() == ast.EOF {
					return true
				}
			} else { // It's an empty result
				rp.result = ast.AppendNode(rp.result, ast.NewEmptyNonTerminalNode(rp.token, pos, rp.interpreter))
			}
		}
	}
	return false
}

func (rp *recursive) parseNext(i int, node parsley.Node, depth int, ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos, mergeCurtailingParsers bool) bool {
	if len(rp.nodes) < depth+1 {
		rp.nodes = append(rp.nodes, node)
	} else {
		rp.nodes[depth] = node
	}
	if i > 0 || node.ReaderPos() > pos {
		leftRecCtx = data.EmptyIntMap
		mergeCurtailingParsers = false
	}
	if rp.parse(depth+1, ctx, leftRecCtx, node.ReaderPos(), mergeCurtailingParsers) {
		return true
	}
	return false
}
