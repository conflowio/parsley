// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package interpreter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/ast/interpreter"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/parsley/parsleyfakes"
)

var _ = Describe("Interpreter", func() {
	var (
		node         *parsleyfakes.FakeNonTerminalNode
		node1, node2 *parsleyfakes.FakeNonLiteralNode
		ctx          interface{}
	)

	BeforeEach(func() {
		ctx = "context"
		node1 = &parsleyfakes.FakeNonLiteralNode{}
		node1.ValueReturns(1, parsley.NewErrorf(parsley.Pos(1), "err1"))
		node1.SchemaReturns("testtype")
		node2 = &parsleyfakes.FakeNonLiteralNode{}
		node2.ValueReturns(2, parsley.NewErrorf(parsley.Pos(2), "err2"))
		node = &parsleyfakes.FakeNonTerminalNode{}
		node.ChildrenReturns([]parsley.Node{node1, node2})
	})

	Describe("Select", func() {
		It("should return the value of the selected node", func() {
			f := interpreter.Select(1)

			val, err := f.Eval(ctx, node)
			Expect(val).To(Equal(2))
			Expect(err).To(MatchError("err2"))

			Expect(node1.ValueCallCount()).To(Equal(0))
			Expect(node2.ValueCallCount()).To(Equal(1))
			passedCtx := node2.ValueArgsForCall(0)
			Expect(passedCtx).To(Equal(ctx))
		})

		Context("when the node index is out of bounds (+)", func() {
			It("should panic", func() {
				f := interpreter.Select(2)
				Expect(func() { f.Eval(ctx, node) }).To(Panic())
			})
		})

		Context("when the node index is out of bounds (-)", func() {
			It("should panic", func() {
				f := interpreter.Select(-1)
				Expect(func() { f.Eval(ctx, node) }).To(Panic())
			})
		})

		Context("static checking", func() {
			Context("when the node is not static checkable", func() {
				It("should return with the type of the indexed node", func() {
					f := interpreter.Select(0)
					nodeType, err := f.(parsley.StaticChecker).StaticCheck(ctx, node)
					Expect(err).ToNot(HaveOccurred())
					Expect(nodeType).To(Equal("testtype"))
				})
			})

			Context("when the node index is out of bounds (+)", func() {
				It("should panic", func() {
					f := interpreter.Select(2)
					Expect(func() { f.(parsley.StaticChecker).StaticCheck(ctx, node) }).To(Panic())
				})
			})

			Context("when the node index is out of bounds (-)", func() {
				It("should panic", func() {
					f := interpreter.Select(-1)
					Expect(func() { f.(parsley.StaticChecker).StaticCheck(ctx, node) }).To(Panic())
				})
			})
		})
	})

	Describe("Nil", func() {
		It("should always return nil", func() {
			ctx := "context"
			f := interpreter.Nil()
			node1 := &parsleyfakes.FakeNonLiteralNode{}
			node1.ValueReturns(1, parsley.NewErrorf(parsley.Pos(1), "err1"))
			node := &parsleyfakes.FakeNonTerminalNode{}
			node.ChildrenReturns([]parsley.Node{node1})
			val, err := f.Eval(ctx, node)
			Expect(val).To(BeNil())
			Expect(err).ToNot(HaveOccurred())
			Expect(node1.ValueCallCount()).To(Equal(0))
		})
	})

	Describe("Array", func() {
		var (
			node                   *parsleyfakes.FakeNonTerminalNode
			child1, child2, child3 *parsleyfakes.FakeNonLiteralNode
			value                  interface{}
			evalErr                parsley.Error
		)

		BeforeEach(func() {
			child1 = &parsleyfakes.FakeNonLiteralNode{}
			child1.ValueReturns("v1", nil)
			child2 = &parsleyfakes.FakeNonLiteralNode{}
			child2.ValueReturns(",", nil)
			child3 = &parsleyfakes.FakeNonLiteralNode{}
			child3.ValueReturns("v3", nil)
			node = &parsleyfakes.FakeNonTerminalNode{}
			node.ChildrenReturns([]parsley.Node{child1, child2, child3})
		})

		JustBeforeEach(func() {
			value, evalErr = interpreter.Array().Eval(ctx, node)
		})

		It("should return with an array", func() {
			Expect(value).To(Equal([]interface{}{"v1", "v3"}))
			Expect(evalErr).ToNot(HaveOccurred())
		})

		Context("when there are no nodes", func() {
			BeforeEach(func() {
				node.ChildrenReturns([]parsley.Node{})
			})
			It("should return with an empty array", func() {
				Expect(value).To(Equal([]interface{}{}))
				Expect(evalErr).ToNot(HaveOccurred())
			})
		})

		Context("when there are even nodes (trailing separator)", func() {
			BeforeEach(func() {
				node.ChildrenReturns([]parsley.Node{child1, child2})
			})
			It("should return with an array", func() {
				Expect(value).To(Equal([]interface{}{"v1"}))
				Expect(evalErr).ToNot(HaveOccurred())
			})
		})

		Context("when a node evaluation has an error", func() {
			var err = parsley.NewErrorf(parsley.Pos(1), "some error")
			BeforeEach(func() {
				child1.ValueReturns(nil, err)
			})
			It("returns with the error", func() {
				Expect(value).To(BeNil())
				Expect(evalErr).To(MatchError(err))
			})
		})
	})

	Describe("Object", func() {
		var (
			node1, node3 parsley.NonTerminalNode
			node2        *parsleyfakes.FakeNode
			node         *parsleyfakes.FakeNonTerminalNode
			value        interface{}
			evalErr      parsley.Error
		)

		BeforeEach(func() {
			c1 := &parsleyfakes.FakeNonLiteralNode{}
			c1.ValueReturns("key1", nil)
			c2 := &parsleyfakes.FakeNonLiteralNode{}
			c2.ValueReturns(":", nil)
			c3 := &parsleyfakes.FakeNonLiteralNode{}
			c3.ValueReturns("value1", nil)
			node1 = ast.NewNonTerminalNode("KEY_VALUE", []parsley.Node{c1, c2, c3}, nil)

			node2 = &parsleyfakes.FakeNode{}

			c4 := &parsleyfakes.FakeNonLiteralNode{}
			c4.ValueReturns("key2", nil)
			c5 := &parsleyfakes.FakeNonLiteralNode{}
			c5.ValueReturns(":", nil)
			c6 := &parsleyfakes.FakeNonLiteralNode{}
			c6.ValueReturns("value2", nil)
			node3 = ast.NewNonTerminalNode("KEY_VALUE", []parsley.Node{c4, c5, c6}, nil)

			node = &parsleyfakes.FakeNonTerminalNode{}
			node.ChildrenReturns([]parsley.Node{node1, node2, node3})
		})

		JustBeforeEach(func() {
			value, evalErr = interpreter.Object().Eval(ctx, node)
		})

		It("should return with an object", func() {
			Expect(value).To(Equal(map[string]interface{}{"key1": "value1", "key2": "value2"}))
			Expect(evalErr).ToNot(HaveOccurred())
		})

		Context("when there are no nodes", func() {
			BeforeEach(func() {
				node.ChildrenReturns([]parsley.Node{})
			})
			It("should return with an empty map", func() {
				Expect(value).To(Equal(map[string]interface{}{}))
				Expect(evalErr).ToNot(HaveOccurred())
			})
		})

		Context("when there are even nodes (trailing separator)", func() {
			BeforeEach(func() {
				node.ChildrenReturns([]parsley.Node{node1, node2})
			})
			It("should return with a map", func() {
				Expect(value).To(Equal(map[string]interface{}{"key1": "value1"}))
				Expect(evalErr).ToNot(HaveOccurred())
			})
		})

		Context("when a key node evaluation has an error", func() {
			var err = parsley.NewErrorf(parsley.Pos(1), "some error")
			BeforeEach(func() {
				node1.Children()[0].(*parsleyfakes.FakeNonLiteralNode).ValueReturns(0, err)
			})
			It("returns with the error", func() {
				Expect(value).To(BeNil())
				Expect(evalErr).To(MatchError(err))
			})
		})

		Context("when a value node evaluation has an error", func() {
			var err = parsley.NewErrorf(parsley.Pos(1), "some error")
			BeforeEach(func() {
				node1.Children()[2].(*parsleyfakes.FakeNonLiteralNode).ValueReturns(0, err)
			})
			It("returns with the error", func() {
				Expect(value).To(BeNil())
				Expect(evalErr).To(MatchError(err))
			})
		})
	})
})
