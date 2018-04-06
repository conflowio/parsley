// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("NodeList", func() {

	var (
		nl, nl2        ast.NodeList
		n1, n2, n3, n4 *parsleyfakes.FakeNode
	)

	BeforeEach(func() {
		n1 = &parsleyfakes.FakeNode{}
		n1.TokenReturns("n1")
		n2 = &parsleyfakes.FakeNode{}
		n2.TokenReturns("n2")
		n3 = &parsleyfakes.FakeNode{}
		n3.TokenReturns("n3")
		n4 = &parsleyfakes.FakeNode{}
		n4.TokenReturns("n4")

		nl2 = ast.NodeList([]parsley.Node{n3, n4})
	})

	It("should have a non-empty token", func() {
		Expect(nl.Token()).ToNot(BeEmpty())
	})

	Describe("Value", func() {
		Context("when empty", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{})
			})
			It("should return nil", func() {
				Expect(nl.Value(nil)).To(BeNil())
			})
		})

		Context("when not empty", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{n1, n2})
			})
			It("should return the value of the first item", func() {
				err := parsley.NewError(parsley.Pos(1), "some error")
				n1.ValueReturns("res", err)

				ctx := "foo"
				val, err := nl.Value(ctx)
				Expect(val).To(Equal("res"))
				Expect(err).To(MatchError(err))

				Expect(n1.ValueCallCount()).To(Equal(1))
				passedCtx := n1.ValueArgsForCall(0)
				Expect(passedCtx).To(Equal(ctx))

				Expect(n2.ValueCallCount()).To(Equal(0))
			})
		})
	})

	Describe("Pos", func() {
		Context("when empty", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{})
			})
			It("should return nil position", func() {
				Expect(nl.Pos()).To(Equal(parsley.NilPos))
			})
		})

		Context("when not empty", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{n1, n2})
				n1.PosReturns(parsley.Pos(1))
				n2.PosReturns(parsley.Pos(2))
			})
			It("should return the position of the first item", func() {
				Expect(nl.Pos()).To(Equal(n1.Pos()))
			})
		})
	})

	Describe("Append", func() {
		Context("when empty", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{})
			})

			It("should append a new item", func() {
				nl.Append(n1)
				Expect(nl).To(Equal(ast.NodeList([]parsley.Node{n1})))
			})

			It("should append a new item list", func() {
				nl.Append(nl2)
				Expect(nl).To(Equal(ast.NodeList([]parsley.Node{n3, n4})))
			})
		})

		Context("when already has items", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{n1})
			})

			It("should append a new item", func() {
				nl.Append(n2)
				Expect(nl).To(Equal(ast.NodeList([]parsley.Node{n1, n2})))
			})

			It("should append a new item list", func() {
				nl.Append(nl2)
				Expect(nl).To(Equal(ast.NodeList([]parsley.Node{n1, n3, n4})))
			})

			It("should only append an empty node once", func() {
				nl.Append(ast.EmptyNode(1))
				nl.Append(ast.EmptyNode(1))
				Expect(nl).To(Equal(ast.NodeList([]parsley.Node{n1, ast.EmptyNode(1)})))
			})
		})
	})

	Describe("ReaderPos", func() {
		It("should panic", func() {
			Expect(func() { nl.ReaderPos() }).To(Panic())
		})
	})

})

var _ = Describe("AppendNode", func() {

	var (
		n1, n2, n3, n4 *parsleyfakes.FakeNode
	)

	BeforeEach(func() {
		n1 = &parsleyfakes.FakeNode{}
		n1.TokenReturns("n1")
		n2 = &parsleyfakes.FakeNode{}
		n2.TokenReturns("n2")
		n3 = &parsleyfakes.FakeNode{}
		n3.TokenReturns("n3")
		n4 = &parsleyfakes.FakeNode{}
		n4.TokenReturns("n4")
	})

	DescribeTable("result",
		func(n1 parsley.Node, n2 parsley.Node, res parsley.Node) {
			if res != nil {
				Expect(ast.AppendNode(n1, n2)).To(Equal(res))
			} else {
				Expect(ast.AppendNode(n1, n2)).To(BeNil())
			}
		},
		Entry("nil + nil", nil, nil, nil),
		Entry("node + nil", n1, nil, n1),
		Entry("nil + node", nil, n1, n1),
		Entry("node + node", n1, n2, ast.NodeList([]parsley.Node{n1, n2})),
		Entry("node + list", n1, ast.NodeList([]parsley.Node{n2, n3}), ast.NodeList([]parsley.Node{n1, n2, n3})),
		Entry("list + node", ast.NodeList([]parsley.Node{n1, n2}), n3, ast.NodeList([]parsley.Node{n1, n2, n3})),
		Entry(
			"list + list",
			ast.NodeList([]parsley.Node{n1, n2}),
			ast.NodeList([]parsley.Node{n3, n4}),
			ast.NodeList([]parsley.Node{n1, n2, n3, n4}),
		),
	)
})
