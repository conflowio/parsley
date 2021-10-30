// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/ast/astfakes"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/parsley/parsleyfakes"
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

	Describe("Token", func() {
		Context("when empty", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{})
			})
			It("should panic", func() {
				Expect(func() { nl.Token() }).To(Panic())
			})
		})

		Context("when not empty", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{n1, n2})
			})
			It("should return the reader pos of the first item", func() {
				n1.TokenReturns("testtoken")

				Expect(nl.Token()).To(Equal("testtoken"))

				Expect(n1.TokenCallCount()).To(Equal(1))
				Expect(n2.TokenCallCount()).To(Equal(0))
			})
		})
	})

	Describe("Pos", func() {
		Context("when empty", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{})
			})
			It("should panic", func() {
				Expect(func() { nl.Pos() }).To(Panic())
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
		Context("when empty", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{})
			})
			It("should panic", func() {
				Expect(func() { nl.ReaderPos() }).To(Panic())
			})
		})

		Context("when not empty", func() {
			BeforeEach(func() {
				nl = ast.NodeList([]parsley.Node{n1, n2})
			})
			It("should return the reader pos of the first item", func() {
				n1.ReaderPosReturns(parsley.Pos(2))

				Expect(nl.ReaderPos()).To(Equal(parsley.Pos(2)))

				Expect(n1.ReaderPosCallCount()).To(Equal(1))
				Expect(n2.ReaderPosCallCount()).To(Equal(0))
			})
		})
	})

	Describe("Walk", func() {
		It("should call the function on the first node", func() {
			nl = ast.NodeList([]parsley.Node{n1, n2})
			called := []parsley.Node{}
			f := func(n parsley.Node) bool {
				called = append(called, n)
				return false
			}

			res := nl.Walk(f)
			Expect(res).To(BeFalse())

			Expect(called).To(Equal([]parsley.Node{n1}))
		})
	})

	Describe("SetReaderPos", func() {
		It("should call the function on all nodes", func() {
			n1 := &astfakes.FakeReaderPosSetterNode{}
			n2 := &astfakes.FakeReaderPosSetterNode{}
			nl = ast.NodeList([]parsley.Node{n1, n2})
			f := func(pos parsley.Pos) parsley.Pos { return pos }

			nl.SetReaderPos(f)

			Expect(n1.SetReaderPosCallCount()).To(Equal(1))
			Expect(n2.SetReaderPosCallCount()).To(Equal(1))
		})
	})
})
