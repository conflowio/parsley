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
	var cp, qcp data.IntSet

	JustBeforeEach(func() {
		q = &parsleyfakes.FakeParser{}
		q.ParseReturnsOnCall(0, qres, qcp)
		p := combinator.Single(q)
		res, cp = p.Parse(ctx, leftRecCtx, pos)
	})

	BeforeEach(func() {
		ctx = parsley.NewContext(&parsleyfakes.FakeReader{})
		leftRecCtx = data.NewIntMap(map[int]int{1: 2})
		pos = parsley.Pos(1)
		qcp = data.NewIntSet(1, 2)
	})

	Context("when q returns with a non terminal node with a single child", func() {
		var childNode parsley.Node

		BeforeEach(func() {
			childNode = ast.NewTerminalNode("child", "child", parsley.Pos(1), parsley.Pos(2))
			qres = ast.NewNonTerminalNode("res", []parsley.Node{childNode}, nil)
		})

		It("should return the single node instead", func() {
			Expect(res).To(BeEquivalentTo(childNode))
			Expect(cp).To(BeEquivalentTo(qcp))
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
		})
	})

	Context("when q does not return with a non terminal node", func() {

		BeforeEach(func() {
			qres = ast.NewTerminalNode("child", "child", parsley.Pos(1), parsley.Pos(2))
		})

		It("should return the original result", func() {
			Expect(res).To(BeEquivalentTo(qres))
			Expect(cp).To(BeEquivalentTo(qcp))
		})
	})

	Context("when q returns nil", func() {

		BeforeEach(func() {
			qres = nil
		})

		It("should return tnil", func() {
			Expect(res).To(BeNil())
			Expect(cp).To(BeEquivalentTo(qcp))
		})
	})

})
