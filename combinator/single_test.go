// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("Single", func() {

	var q *parsleyfakes.FakeParser
	var ctx *parsley.Context
	var leftRecCtx data.IntMap
	var pos parsley.Pos
	var res, qres parsley.Node
	var err, qerr parsley.Error
	var cp, qcp data.IntSet

	JustBeforeEach(func() {
		q = &parsleyfakes.FakeParser{}
		q.ParseReturnsOnCall(0, qres, qcp, qerr)
		p := combinator.Single(q)
		res, cp, err = p.Parse(ctx, leftRecCtx, pos)
	})

	BeforeEach(func() {
		ctx = parsley.NewContext(parsley.NewFileSet(), &parsleyfakes.FakeReader{})
		leftRecCtx = data.NewIntMap(map[int]int{1: 2})
		pos = parsley.Pos(1)
		qcp = data.NewIntSet(1, 2)
	})

	Context("when q returns with a non terminal node with a single child", func() {
		var childNode parsley.Node

		BeforeEach(func() {
			childNode = ast.NewTerminalNode("string", "child", "child", parsley.Pos(1), parsley.Pos(2))
			qres = ast.NewNonTerminalNode("res", []parsley.Node{childNode}, nil)
		})

		It("should return the single node instead", func() {
			Expect(res).To(BeEquivalentTo(childNode))
			Expect(cp).To(BeEquivalentTo(qcp))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when q returns with a non terminal node with multiple children", func() {

		BeforeEach(func() {
			child1 := &parsleyfakes.FakeNode{}
			child2 := &parsleyfakes.FakeNode{}
			qres = ast.NewNonTerminalNode("res", []parsley.Node{child1, child2}, nil)
		})

		It("should return the original result", func() {
			Expect(res).To(BeEquivalentTo(qres))
			Expect(cp).To(BeEquivalentTo(qcp))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when q does not return with a non terminal node", func() {

		BeforeEach(func() {
			qres = ast.NewTerminalNode("string", "child", "child", parsley.Pos(1), parsley.Pos(2))
		})

		It("should return the original result", func() {
			Expect(res).To(BeEquivalentTo(qres))
			Expect(cp).To(BeEquivalentTo(qcp))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when q returns an error", func() {

		BeforeEach(func() {
			qerr = parsley.NewErrorf(parsley.Pos(1), "some error")
		})

		It("should return tnil", func() {
			Expect(res).To(BeNil())
			Expect(cp).To(BeEquivalentTo(qcp))
			Expect(err).To(BeEquivalentTo(qerr))
		})
	})

})
