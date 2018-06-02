// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("TerminalNode", func() {
	var (
		node      *ast.TerminalNode
		token     string      = "TEST"
		value     interface{} = "some value"
		pos       parsley.Pos = parsley.Pos(1)
		readerPos parsley.Pos = parsley.Pos(2)
	)

	JustBeforeEach(func() {
		node = ast.NewTerminalNode(token, value, pos, readerPos)
	})

	Describe("Methods", func() {
		It("Token() should return with the token value", func() {
			Expect(node.Token()).To(Equal(token))
		})

		It("Value() should return with the value", func() {
			nodeValue, err := node.Value(nil)
			Expect(nodeValue).To(Equal(value))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Pos() should return with the token value", func() {
			Expect(node.Pos()).To(Equal(pos))
		})

		It("ReaderPos() should return with the reader position", func() {
			Expect(node.ReaderPos()).To(Equal(readerPos))
		})

		It("SetReaderPos() should modify the reader position", func() {
			node.SetReaderPos(func(pos parsley.Pos) parsley.Pos {
				return parsley.Pos(pos + 1)
			})
			Expect(node.ReaderPos()).To(Equal(parsley.Pos(3)))
		})

		It("String() should return with a readable representation", func() {
			Expect(node.String()).To(Equal("TEST{some value, 1..2}"))
		})
	})
})

var _ = Describe("NilNode", func() {
	var (
		node ast.NilNode
		pos  parsley.Pos = parsley.Pos(1)
	)

	JustBeforeEach(func() {
		node = ast.NilNode(pos)
	})

	Describe("Methods", func() {
		It("Token() should return with the token value", func() {
			Expect(node.Token()).To(Equal(ast.NIL))
		})

		It("Value() should return with nil", func() {
			nodeValue, err := node.Value(nil)
			Expect(nodeValue).To(BeNil())
			Expect(err).ToNot(HaveOccurred())
		})

		It("Pos() should return with the token value", func() {
			Expect(node.Pos()).To(Equal(pos))
		})

		It("ReaderPos() should return with the reader position", func() {
			Expect(node.ReaderPos()).To(Equal(pos))
		})

		It("String() should return with NIL", func() {
			Expect(node.String()).To(Equal(ast.NIL))
		})
	})
})

var _ = Describe("NonTerminalNode", func() {
	var (
		node            *ast.NonTerminalNode
		token           string = "TEST"
		children        []parsley.Node
		child1, child2  *parsleyfakes.FakeNode
		pos             parsley.Pos = parsley.Pos(1)
		readerPos       parsley.Pos = parsley.Pos(2)
		interpreter     parsley.Interpreter
		fakeInterpreter *parsleyfakes.FakeInterpreter
		value           = "someValue"
		evalErr         = parsley.NewErrorf(parsley.Pos(1), "eval error")
	)

	BeforeEach(func() {
		fakeInterpreter = &parsleyfakes.FakeInterpreter{}
		fakeInterpreter.EvalReturns(value, evalErr)
		interpreter = fakeInterpreter
		child1 = &parsleyfakes.FakeNode{}
		child1.PosReturns(pos)
		child2 = &parsleyfakes.FakeNode{}
		child2.ReaderPosReturns(readerPos)
		children = []parsley.Node{child1, child2}
	})

	Context("when NewNonTerminalNode is called with no children", func() {
		It("should panic", func() {
			Expect(func() { ast.NewNonTerminalNode(token, nil, interpreter) }).To(Panic())
		})
	})

	Context("when NewNonTerminalNode is called with nil children", func() {
		It("should panic", func() {
			Expect(func() { ast.NewNonTerminalNode(token, []parsley.Node{nil}, interpreter) }).To(Panic())
		})
	})

	Context("when created with child nodes", func() {
		JustBeforeEach(func() {
			node = ast.NewNonTerminalNode(token, children, interpreter)
		})

		It("should get the position from the first node", func() {
			Expect(child1.PosCallCount()).To(Equal(1))
			Expect(child2.PosCallCount()).To(Equal(0))
		})

		It("should get the reader positon from the last node", func() {
			Expect(child1.ReaderPosCallCount()).To(Equal(0))
			Expect(child2.ReaderPosCallCount()).To(Equal(1))
		})

		Describe("Methods", func() {
			It("Token() should return with the token value", func() {
				Expect(node.Token()).To(Equal(token))
			})

			It("Value() should return the result of the interpreter", func() {
				ctx := "some context"
				nodeValue, err := node.Value(ctx)
				Expect(nodeValue).To(Equal(value))
				Expect(err).To(MatchError(evalErr))

				Expect(fakeInterpreter.EvalCallCount()).To(Equal(1))
				passedCtx, passedNodes := fakeInterpreter.EvalArgsForCall(0)
				Expect(passedCtx).To(Equal(ctx))
				Expect(passedNodes).To(Equal(children))
			})

			Context("when there is no interpreter", func() {
				BeforeEach(func() {
					interpreter = nil
				})
				It("Value() should panic", func() {
					Expect(func() { node.Value("ctx") }).To(Panic())
				})
			})

			It("Pos() should return with the token value", func() {
				Expect(node.Pos()).To(Equal(pos))
			})

			It("ReaderPos() should return with the reader position", func() {
				Expect(node.ReaderPos()).To(Equal(readerPos))
			})

			It("SetReaderPos() should modify the reader position", func() {
				node.SetReaderPos(func(pos parsley.Pos) parsley.Pos {
					return parsley.Pos(pos + 1)
				})
				Expect(node.ReaderPos()).To(Equal(parsley.Pos(3)))
			})

			It("Children() should return with the children", func() {
				Expect(node.Children()).To(Equal(children))
			})

			Context("when having real children", func() {
				BeforeEach(func() {
					children = []parsley.Node{
						ast.NewTerminalNode("STRING", "foo", parsley.Pos(1), parsley.Pos(2)),
					}
				})
				It("String() should return with a readable representation", func() {
					Expect(node.String()).To(Equal("TEST{[STRING{foo, 1..2}], 1..2}"))
				})
			})
		})

	})

	Context("when created without child nodes", func() {
		JustBeforeEach(func() {
			node = ast.NewEmptyNonTerminalNode(token, pos, interpreter)
		})

		Describe("Methods", func() {
			It("Token() should return with the token value", func() {
				Expect(node.Token()).To(Equal(token))
			})

			It("Value() should return the result of the interpreter", func() {
				ctx := "some context"
				nodeValue, err := node.Value(ctx)
				Expect(nodeValue).To(Equal(value))
				Expect(err).To(MatchError(evalErr))

				Expect(fakeInterpreter.EvalCallCount()).To(Equal(1))
				passedCtx, passedNodes := fakeInterpreter.EvalArgsForCall(0)
				Expect(passedCtx).To(Equal(ctx))
				Expect(passedNodes).To(Equal([]parsley.Node{}))
			})

			Context("when there is no interpreter", func() {
				BeforeEach(func() {
					interpreter = nil
				})
				It("Value() should panic", func() {
					Expect(func() { node.Value("ctx") }).To(Panic())
				})
			})

			It("Pos() should return with the token value", func() {
				Expect(node.Pos()).To(Equal(pos))
			})

			It("ReaderPos() should return with the reader position", func() {
				Expect(node.ReaderPos()).To(Equal(pos))
			})

			It("SetReaderPos() should modify the reader position", func() {
				node.SetReaderPos(func(pos parsley.Pos) parsley.Pos {
					return parsley.Pos(pos + 1)
				})
				Expect(node.ReaderPos()).To(Equal(parsley.Pos(2)))
			})

			It("Children() should return nil", func() {
				Expect(node.Children()).To(Equal([]parsley.Node{}))
			})
		})
	})

})

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
				err := parsley.NewErrorf(parsley.Pos(1), "some error")
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
				nl.Append(ast.NilNode(1))
				nl.Append(ast.NilNode(1))
				Expect(nl).To(Equal(ast.NodeList([]parsley.Node{n1, ast.NilNode(1)})))
			})
		})
	})

	Describe("ReaderPos", func() {
		It("should panic", func() {
			Expect(func() { nl.ReaderPos() }).To(Panic())
		})
	})

	Describe("Walk", func() {
		It("should call the function with all children", func() {
			nl = ast.NodeList([]parsley.Node{n1, n2})
			called := []parsley.Node{}
			f := func(i int, n parsley.Node) bool {
				Expect(i).To(Equal(len(called)))
				called = append(called, n)
				return false
			}

			nl.Walk(f)

			Expect(called).To(Equal([]parsley.Node{n1, n2}))
		})

		It("should stop if the function returns with true", func() {
			nl = ast.NodeList([]parsley.Node{n1, n2, n3})
			called := []parsley.Node{}
			f := func(i int, n parsley.Node) bool {
				Expect(i).To(Equal(len(called)))
				called = append(called, n)
				return i > 0
			}

			nl.Walk(f)

			Expect(called).To(Equal([]parsley.Node{n1, n2}))
		})
	})
})
