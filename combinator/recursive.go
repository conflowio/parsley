// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator

import (
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/reader"
)

// recursive is a recursive and-type combinator
type recursive struct {
	nodeBuilder       ast.NodeBuilder
	parserLookUp      func(i int) parser.Parser
	min               int
	max               int
	curtailingParsers data.IntSet
	resultSet         parser.ResultSet
	err               reader.Error
	nodes             []ast.Node
}

// newRecursive creates a new recursive combinator
func newRecursive(nodeBuilder ast.NodeBuilder, parserLookUp func(i int) parser.Parser, min int, max int) *recursive {
	if nodeBuilder == nil {
		panic("Node builder can not be nil!")
	}
	return &recursive{
		nodeBuilder:       nodeBuilder,
		parserLookUp:      parserLookUp,
		min:               min,
		max:               max,
		curtailingParsers: parser.NoCurtailingParsers(),
		resultSet:         parser.ResultSet{},
		nodes:             []ast.Node{},
	}
}

// Parse runs the recursive parser
func (rp *recursive) Parse(h *parser.History, leftRecCtx data.IntMap, r reader.Reader) (data.IntSet, parser.ResultSet, reader.Error) {
	rp.parse(0, h, leftRecCtx, r, true)
	return rp.curtailingParsers, rp.resultSet, rp.err
}

func (rp *recursive) parse(depth int, h *parser.History, leftRecCtx data.IntMap, r reader.Reader, mergeCurtailingParsers bool) bool {
	var cp data.IntSet
	var rs parser.ResultSet
	var err reader.Error
	nextParser := rp.parserLookUp(depth)
	if nextParser != nil {
		h.RegisterCall()
		cp, rs, err = nextParser.Parse(h, leftRecCtx, r.Clone())
		if err != nil && (rp.err == nil || err.Pos().Pos() >= rp.err.Pos().Pos()) {
			rp.err = err
		}
	}

	if mergeCurtailingParsers {
		rp.curtailingParsers = rp.curtailingParsers.Union(cp)
	}

	if len(rs) > 0 {
		for i, result := range rs {
			if len(rp.nodes) < depth+1 {
				rp.nodes = append(rp.nodes, result.Node())
			} else {
				rp.nodes[depth] = result.Node()
			}
			if i > 0 || result.Reader().Cursor().Pos() > r.Cursor().Pos() {
				leftRecCtx = parser.EmptyLeftRecCtx()
				mergeCurtailingParsers = false
			}
			if rp.parse(depth+1, h, leftRecCtx, result.Reader().Clone(), mergeCurtailingParsers) {
				return true
			}
		}
	}
	if len(rs) == 0 {
		if depth >= rp.min && (rp.max == -1 || depth <= rp.max) {
			if depth > 0 {
				nodesCopy := make([]ast.Node, depth)
				copy(nodesCopy[0:depth], rp.nodes[0:depth])
				rp.resultSet.Append(parser.NewResult(rp.nodeBuilder.BuildNode(nodesCopy), r))
				if rp.nodes[depth-1] != nil && rp.nodes[depth-1].Token() == ast.EOF {
					return true
				}
			} else { // It's an empty result
				rp.resultSet.Append(parser.NewResult(rp.nodeBuilder.BuildNode(nil), r))
			}
		}
	}
	return false
}
