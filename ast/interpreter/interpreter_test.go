package interpreter_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/ast/interpreter"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("Interpreter", func() {
	var (
		node1, node2 *parsleyfakes.FakeNode
		ctx          interface{}
	)

	BeforeEach(func() {
		node1 = &parsleyfakes.FakeNode{}
		node1.ValueReturns(1, parsley.NewError(parsley.Pos(1), "err1"))
		node2 = &parsleyfakes.FakeNode{}
		node2.ValueReturns(2, parsley.NewError(parsley.Pos(2), "err2"))
		ctx = "context"
	})

	Describe("Select", func() {
		It("should return the value of the selected node", func() {
			f := interpreter.Select(1)

			val, err := f.Eval(ctx, []parsley.Node{node1, node2})
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
				Expect(func() { f.Eval(ctx, []parsley.Node{node1, node2}) }).To(Panic())
			})
		})

		Context("when the node index is out of bounds (-)", func() {
			It("should panic", func() {
				f := interpreter.Select(-1)
				Expect(func() { f.Eval(ctx, []parsley.Node{node1, node2}) }).To(Panic())
			})
		})
	})

	Describe("Nil", func() {
		It("should always return nil", func() {
			ctx := "context"
			f := interpreter.Nil()
			node1 := &parsleyfakes.FakeNode{}
			node1.ValueReturns(1, parsley.NewError(parsley.Pos(1), "err1"))
			val, err := f.Eval(ctx, []parsley.Node{node1})
			Expect(val).To(BeNil())
			Expect(err).ToNot(HaveOccurred())
			Expect(node1.ValueCallCount()).To(Equal(0))
		})
	})

	Describe("Array", func() {
		var (
			child1, child2, child3 *parsleyfakes.FakeNode
			nodes                  []parsley.Node
			value                  interface{}
			evalErr                parsley.Error
		)

		BeforeEach(func() {
			child1 = &parsleyfakes.FakeNode{}
			child1.ValueReturns("v1", nil)
			child2 = &parsleyfakes.FakeNode{}
			child2.ValueReturns(",", nil)
			child3 = &parsleyfakes.FakeNode{}
			child3.ValueReturns("v3", nil)
			nodes = []parsley.Node{child1, child2, child3}
		})

		JustBeforeEach(func() {
			value, evalErr = interpreter.Array().Eval(ctx, nodes)
		})

		It("should return with an array", func() {
			Expect(value).To(Equal([]interface{}{"v1", "v3"}))
			Expect(evalErr).ToNot(HaveOccurred())
		})

		Context("when there are no nodes", func() {
			BeforeEach(func() {
				nodes = []parsley.Node{}
			})
			It("should return with an empty array", func() {
				Expect(value).To(Equal([]interface{}{}))
				Expect(evalErr).ToNot(HaveOccurred())
			})
		})

		Context("when there are even nodes (trailing separator)", func() {
			BeforeEach(func() {
				nodes = []parsley.Node{child1, child2}
			})
			It("should return with an array", func() {
				Expect(value).To(Equal([]interface{}{"v1"}))
				Expect(evalErr).ToNot(HaveOccurred())
			})
		})

		Context("when a node evaluation has an error", func() {
			var err = parsley.NewError(parsley.Pos(1), "some error")
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
			node1, node3 *ast.NonTerminalNode
			node2        *parsleyfakes.FakeNode
			nodes        []parsley.Node
			value        interface{}
			evalErr      parsley.Error
		)

		BeforeEach(func() {
			c1 := &parsleyfakes.FakeNode{}
			c1.ValueReturns("key1", nil)
			c2 := &parsleyfakes.FakeNode{}
			c2.ValueReturns(":", nil)
			c3 := &parsleyfakes.FakeNode{}
			c3.ValueReturns("value1", nil)
			node1 = ast.NewNonTerminalNode("KEY_VALUE", []parsley.Node{c1, c2, c3}, nil)

			node2 = &parsleyfakes.FakeNode{}

			c4 := &parsleyfakes.FakeNode{}
			c4.ValueReturns("key2", nil)
			c5 := &parsleyfakes.FakeNode{}
			c5.ValueReturns(":", nil)
			c6 := &parsleyfakes.FakeNode{}
			c6.ValueReturns("value2", nil)
			node3 = ast.NewNonTerminalNode("KEY_VALUE", []parsley.Node{c4, c5, c6}, nil)

			nodes = []parsley.Node{node1, node2, node3}
		})

		JustBeforeEach(func() {
			value, evalErr = interpreter.Object().Eval(ctx, nodes)
		})

		It("should return with an object", func() {
			Expect(value).To(Equal(map[string]interface{}{"key1": "value1", "key2": "value2"}))
			Expect(evalErr).ToNot(HaveOccurred())
		})

		Context("when there are no nodes", func() {
			BeforeEach(func() {
				nodes = []parsley.Node{}
			})
			It("should return with an empty map", func() {
				Expect(value).To(Equal(map[string]interface{}{}))
				Expect(evalErr).ToNot(HaveOccurred())
			})
		})

		Context("when there are even nodes (trailing separator)", func() {
			BeforeEach(func() {
				nodes = []parsley.Node{node1, node2}
			})
			It("should return with a map", func() {
				Expect(value).To(Equal(map[string]interface{}{"key1": "value1"}))
				Expect(evalErr).ToNot(HaveOccurred())
			})
		})

		Context("when a key node evaluation has an error", func() {
			var err = parsley.NewError(parsley.Pos(1), "some error")
			BeforeEach(func() {
				nodes[0].(*ast.NonTerminalNode).Children()[0].(*parsleyfakes.FakeNode).ValueReturns(0, err)
			})
			It("returns with the error", func() {
				Expect(value).To(BeNil())
				Expect(evalErr).To(MatchError(err))
			})
		})

		Context("when a value node evaluation has an error", func() {
			var err = parsley.NewError(parsley.Pos(1), "some error")
			BeforeEach(func() {
				nodes[0].(*ast.NonTerminalNode).Children()[2].(*parsleyfakes.FakeNode).ValueReturns(0, err)
			})
			It("returns with the error", func() {
				Expect(value).To(BeNil())
				Expect(evalErr).To(MatchError(err))
			})
		})
	})
})
