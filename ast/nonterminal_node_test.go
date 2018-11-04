// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("NonTerminalNode", func() {
	var (
		node            *ast.NonTerminalNode
		token           string = "TEST"
		children        []parsley.Node
		child1, child2  *parsleyfakes.FakeStaticCheckableNode
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
		child1 = &parsleyfakes.FakeStaticCheckableNode{}
		child1.PosReturns(pos)
		child2 = &parsleyfakes.FakeStaticCheckableNode{}
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
				passedCtx, passedNode := fakeInterpreter.EvalArgsForCall(0)
				Expect(passedCtx).To(Equal(ctx))
				Expect(passedNode).To(Equal(node))
			})

			Context("when there is no interpreter", func() {
				BeforeEach(func() {
					interpreter = nil
				})
				It("Value() should panic", func() {
					Expect(func() { node.Value("ctx") }).To(Panic())
				})

				It("StaticCheck() should panic", func() {
					Expect(func() { node.StaticCheck("ctx") }).To(Panic())
				})
			})

			It("Pos() should return with the token value", func() {
				Expect(node.Pos()).To(Equal(pos))
			})

			It("Type() should return empty before running a static check", func() {
				Expect(node.Type()).To(BeEmpty())
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
						ast.NewTerminalNode("STRING", "foo", "string", parsley.Pos(1), parsley.Pos(2)),
					}
				})
				It("String() should return with a readable representation", func() {
					Expect(node.String()).To(Equal("TEST{[STRING{<string> foo, 1..2}], 1..2}"))
				})
				Context("having a type", func() {
					BeforeEach(func() {
						staticCheckInterpreter := &parsleyfakes.FakeStaticCheckerInterpreter{}
						staticCheckInterpreter.StaticCheckReturns("testtype", nil)
						interpreter = staticCheckInterpreter
					})
					It("String() should return with a readable representation", func() {
						node.StaticCheck("ctx")
						Expect(node.String()).To(Equal("TEST{<testtype> [STRING{<string> foo, 1..2}], 1..2}"))
					})
				})
			})

			Describe("StaticCheck()", func() {
				var staticCheckInterpreter *parsleyfakes.FakeStaticCheckerInterpreter

				Context("when there is no error", func() {
					BeforeEach(func() {
						staticCheckInterpreter = &parsleyfakes.FakeStaticCheckerInterpreter{}
						staticCheckInterpreter.StaticCheckReturns("testtype", nil)
						interpreter = staticCheckInterpreter
					})

					It("should set the type", func() {
						ctx := "some context"
						err := node.StaticCheck(ctx)
						Expect(err).ToNot(HaveOccurred())

						Expect(staticCheckInterpreter.StaticCheckCallCount()).To(Equal(1))
						passedCtx, passedNode := staticCheckInterpreter.StaticCheckArgsForCall(0)
						Expect(passedCtx).To(Equal(ctx))
						Expect(passedNode).To(Equal(node))

						Expect(node.Type()).To(Equal("testtype"))
					})

					It("should check the children", func() {
						ctx := "some context"
						err := node.StaticCheck(ctx)
						Expect(err).ToNot(HaveOccurred())

						Expect(child1.StaticCheckCallCount()).To(Equal(1))
						Expect(child2.StaticCheckCallCount()).To(Equal(1))

						passedCtx1 := child1.StaticCheckArgsForCall(0)
						Expect(passedCtx1).To(Equal(ctx))

						passedCtx2 := child2.StaticCheckArgsForCall(0)
						Expect(passedCtx2).To(Equal(ctx))
					})
				})

				Context("when there is an error", func() {
					var checkErr parsley.Error

					BeforeEach(func() {
						checkErr = parsley.NewError(parsley.Pos(1), errors.New("check error"))
						staticCheckInterpreter = &parsleyfakes.FakeStaticCheckerInterpreter{}
						staticCheckInterpreter.StaticCheckReturns("", checkErr)
						interpreter = staticCheckInterpreter
					})

					It("should return the error", func() {
						ctx := "some context"
						err := node.StaticCheck(ctx)
						Expect(err).To(MatchError(checkErr))
						Expect(node.Type()).To(BeEmpty())
					})
				})

				Context("when a child has an error", func() {
					var checkErr parsley.Error

					BeforeEach(func() {
						checkErr = parsley.NewError(parsley.Pos(1), errors.New("check error"))
						child1.StaticCheckReturns(checkErr)
					})

					It("should return the error", func() {
						ctx := "some context"
						err := node.StaticCheck(ctx)
						Expect(err).To(MatchError(checkErr))
						Expect(node.Type()).To(BeEmpty())
					})
				})

				Context("when the interpreter is not a static checker", func() {
					It("should return no error", func() {
						ctx := "some context"
						err := node.StaticCheck(ctx)
						Expect(err).ToNot(HaveOccurred())
						Expect(node.Type()).To(BeEmpty())
					})
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
				passedCtx, passedNode := fakeInterpreter.EvalArgsForCall(0)
				Expect(passedCtx).To(Equal(ctx))
				Expect(passedNode).To(Equal(node))
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
