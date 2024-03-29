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

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/parsley/parsleyfakes"
)

var _ = Describe("NonTerminalNode", func() {
	var (
		node            *ast.NonTerminalNode
		token           string = "TEST"
		children        []parsley.Node
		child1, child2  *parsleyfakes.FakeStaticCheckableNode
		child3, child4  *parsleyfakes.FakeTransformableNode
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
		child3 = &parsleyfakes.FakeTransformableNode{}
		child4 = &parsleyfakes.FakeTransformableNode{}
		child4.ReaderPosReturns(readerPos)
		children = []parsley.Node{child1, child2, child3, child4}
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
			Expect(child4.ReaderPosCallCount()).To(Equal(1))
		})

		Describe("Methods", func() {
			It("Token() should return with the token value", func() {
				Expect(node.Token()).To(Equal(token))
			})

			It("EvaluateNode() should return the result of the interpreter", func() {
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
				It("EvaluateNode() should panic", func() {
					Expect(func() { node.Value("ctx") }).To(Panic())
				})

				It("StaticCheck() should return no error", func() {
					Expect(node.StaticCheck("ctx")).ToNot(HaveOccurred())
				})
			})

			It("Pos() should return with the token value", func() {
				Expect(node.Pos()).To(Equal(pos))
			})

			It("Schema() should return nil before running a static check", func() {
				Expect(node.Schema()).To(BeNil())
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
						ast.NewTerminalNode("string", "STRING", "foo", parsley.Pos(1), parsley.Pos(2)),
					}
				})
				It("String() should return with a readable representation", func() {
					Expect(node.String()).To(Equal("TEST{[STRING{foo, 1..2}], 1..2}"))
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

						Expect(node.Schema()).To(Equal("testtype"))
					})
				})

				Context("when there is an error", func() {
					var checkErr parsley.Error

					BeforeEach(func() {
						checkErr = parsley.NewError(parsley.Pos(1), errors.New("check error"))
						staticCheckInterpreter = &parsleyfakes.FakeStaticCheckerInterpreter{}
						staticCheckInterpreter.StaticCheckReturns(nil, checkErr)
						interpreter = staticCheckInterpreter
					})

					It("should return the error", func() {
						ctx := "some context"
						err := node.StaticCheck(ctx)
						Expect(err).To(MatchError(checkErr))
						Expect(node.Schema()).To(BeNil())
					})
				})

				Context("when the interpreter is not a static checker", func() {
					It("should return no error", func() {
						ctx := "some context"
						err := node.StaticCheck(ctx)
						Expect(err).ToNot(HaveOccurred())
						Expect(node.Schema()).To(BeNil())
					})
				})
			})

			Describe("Transform()", func() {
				var transformerInterpreter *parsleyfakes.FakeNodeTransformerInterpreter

				Context("when there is no error", func() {
					transformedNode := &parsleyfakes.FakeNode{}
					transformedNode.TokenReturns("TRANSFORMED")

					BeforeEach(func() {
						transformerInterpreter = &parsleyfakes.FakeNodeTransformerInterpreter{}
						transformerInterpreter.TransformNodeReturns(transformedNode, nil)
						interpreter = transformerInterpreter
					})

					It("should return with the transformed node", func() {
						ctx := "some context"
						res, err := node.Transform(ctx)
						Expect(err).ToNot(HaveOccurred())
						Expect(res).To(Equal(transformedNode))

						Expect(transformerInterpreter.TransformNodeCallCount()).To(Equal(1))
						passedCtx, passedNode := transformerInterpreter.TransformNodeArgsForCall(0)
						Expect(passedCtx).To(Equal(ctx))
						Expect(passedNode).To(Equal(node))
					})
				})

				Context("when there is an error", func() {
					var transformErr parsley.Error

					BeforeEach(func() {
						transformErr = parsley.NewError(parsley.Pos(1), errors.New("transform error"))
						transformerInterpreter = &parsleyfakes.FakeNodeTransformerInterpreter{}
						transformerInterpreter.TransformNodeReturns(nil, transformErr)
						interpreter = transformerInterpreter
					})

					It("should return the error", func() {
						ctx := "some context"
						res, err := node.Transform(ctx)
						Expect(err).To(MatchError(transformErr))
						Expect(res).To(BeNil())
					})
				})

				Context("when the interpreter is not a static checker", func() {
					transformed3 := &parsleyfakes.FakeNode{}
					transformed3.TokenReturns("TR3")
					transformed4 := &parsleyfakes.FakeNode{}
					transformed4.TokenReturns("TR4")

					BeforeEach(func() {
						child3.TransformReturns(transformed3, nil)
						child4.TransformReturns(transformed4, nil)
					})

					It("should return no error", func() {
						ctx := "some context"
						res, err := node.Transform(ctx)

						Expect(err).ToNot(HaveOccurred())
						Expect(res).To(Equal(node))
					})

					It("should call transform on the children", func() {
						ctx := "some context"
						node.Transform(ctx)

						Expect(child3.TransformCallCount()).To(Equal(1))
						Expect(child4.TransformCallCount()).To(Equal(1))

						passedCtx3 := child3.TransformArgsForCall(0)
						Expect(passedCtx3).To(Equal(ctx))
						passedCtx4 := child4.TransformArgsForCall(0)
						Expect(passedCtx4).To(Equal(ctx))

						Expect(node.Children()).To(Equal([]parsley.Node{
							child1, child2, transformed3, transformed4,
						}))
					})

					Context("if a child has a transform error", func() {
						transformErr := parsley.NewErrorf(parsley.Pos(1), "transform error")

						BeforeEach(func() {
							child4.TransformReturns(nil, transformErr)
						})

						It("returns the error", func() {
							res, err := node.Transform("some context")
							Expect(err).To(MatchError(transformErr))
							Expect(res).To(BeNil())
						})
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

			It("EvaluateNode() should return the result of the interpreter", func() {
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
				It("EvaluateNode() should panic", func() {
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
				Expect(node.Children()).To(BeEmpty())
			})
		})
	})

})
