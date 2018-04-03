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

// recursive is a recursive and-type combinator
type recursive struct {
	nodeBuilder       parsley.NodeBuilder
	parserLookUp      func(i int) parsley.Parser
	min               int
	max               int
	curtailingParsers data.IntSet
	result            []parsley.Node
	err               parsley.Error
	nodes             []parsley.Node
}

// newRecursive creates a new recursive combinator
func newRecursive(nodeBuilder parsley.NodeBuilder, parserLookUp func(i int) parsley.Parser, min int, max int) *recursive {
	if nodeBuilder == nil {
		panic("Node builder can not be nil!")
	}
	return &recursive{
		nodeBuilder:       nodeBuilder,
		parserLookUp:      parserLookUp,
		min:               min,
		max:               max,
		curtailingParsers: data.EmptyIntSet(),
		result:            []parsley.Node{},
		nodes:             []parsley.Node{},
	}
}

// Parse runs the recursive parser
func (rp *recursive) Parse(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
	rp.parse(0, h, leftRecCtx, r, pos, true)
	return rp.curtailingParsers, rp.result, rp.err
}

func (rp *recursive) parse(depth int, h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int, mergeCurtailingParsers bool) bool {
	var cp data.IntSet
	var nodes []parsley.Node
	var err parsley.Error
	nextParser := rp.parserLookUp(depth)
	if nextParser != nil {
		h.RegisterCall()
		cp, nodes, err = nextParser.Parse(h, leftRecCtx, r, pos)
		if err != nil && (rp.err == nil || err.Pos() >= rp.err.Pos()) {
			rp.err = err
		}
	}

	if mergeCurtailingParsers {
		rp.curtailingParsers = rp.curtailingParsers.Union(cp)
	}

	if len(nodes) > 0 {
		for i, node := range nodes {
			if len(rp.nodes) < depth+1 {
				rp.nodes = append(rp.nodes, node)
			} else {
				rp.nodes[depth] = node
			}
			if i > 0 || node.ReaderPos() > pos {
				leftRecCtx = data.EmptyIntMap()
				mergeCurtailingParsers = false
			}
			if rp.parse(depth+1, h, leftRecCtx, r, node.ReaderPos(), mergeCurtailingParsers) {
				return true
			}
		}
	}
	if len(nodes) == 0 {
		if depth >= rp.min && (rp.max == -1 || depth <= rp.max) {
			if depth > 0 {
				nodesCopy := make([]parsley.Node, depth)
				copy(nodesCopy[0:depth], rp.nodes[0:depth])
				rp.result = append(rp.result, rp.nodeBuilder.BuildNode(nodesCopy))
				if rp.nodes[depth-1] != nil && rp.nodes[depth-1].Token() == ast.EOF {
					return true
				}
			} else { // It's an empty result
				rp.result = append(rp.result, rp.nodeBuilder.BuildNode(nil))
			}
		}
	}
	return false
}
