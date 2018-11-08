// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("Parse", func() {
	var (
		r         *parsleyfakes.FakeReader
		p         *parsleyfakes.FakeParser
		ctx       *parsley.Context
		userCtx   interface{}
		val       interface{}
		err       error
		parserRes parsley.Node
		parserErr parsley.Error
		res       parsley.Node
	)

	BeforeEach(func() {
		f := &parsleyfakes.FakeFile{}
		f.LenReturns(3)
		position := &parsleyfakes.FakePosition{}
		position.StringReturns("testpos")
		f.PositionReturns(position)

		fs := parsley.NewFileSet(f)

		r = &parsleyfakes.FakeReader{}
		userCtx = "user context"
		ctx = parsley.NewContext(fs, r)
		ctx.SetUserContext(userCtx)
		r.PosReturns(parsley.Pos(1))
		p = &parsleyfakes.FakeParser{}
		resultNode := &parsleyfakes.FakeNode{}
		resultNode.TokenReturns("TEST RESULT")
		parserRes = resultNode
		parserErr = nil
	})

	JustBeforeEach(func() {
		p.ParseReturns(parserRes, data.EmptyIntSet, parserErr)
		res, err = parsley.Parse(ctx, p)
	})

	It("gets the zero position from the reader", func() {
		Expect(r.PosCallCount()).To(Equal(1))
		Expect(r.PosArgsForCall(0)).To(Equal(0))
	})

	It("calls the parser", func() {
		Expect(p.ParseCallCount()).To(Equal(1))
		passedCtx, passedLeftRecCtx, passedPos := p.ParseArgsForCall(0)
		Expect(passedCtx).To(BeEquivalentTo(ctx))
		Expect(passedLeftRecCtx).To(BeEquivalentTo(data.EmptyIntMap))
		Expect(passedPos).To(Equal(parsley.Pos(1)))
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns the result", func() {
		Expect(p.ParseCallCount()).To(Equal(1))
		passedCtx, passedLeftRecCtx, passedPos := p.ParseArgsForCall(0)
		Expect(passedCtx).To(BeEquivalentTo(ctx))
		Expect(passedLeftRecCtx).To(BeEquivalentTo(data.EmptyIntMap))
		Expect(passedPos).To(Equal(parsley.Pos(1)))
		Expect(res).To(BeEquivalentTo(parserRes))
		Expect(err).ToNot(HaveOccurred())
	})

	Context("if transformation is enabled", func() {
		BeforeEach(func() {
			ctx.EnableTransformation()
		})

		Context("if result is transformable", func() {
			var resultNode *parsleyfakes.FakeTransformableNode

			BeforeEach(func() {
				resultNode = &parsleyfakes.FakeTransformableNode{}
				parserRes = resultNode
			})

			It("calls the Transform function", func() {
				Expect(resultNode.TransformCallCount()).To(Equal(1))
				passedUserCtx := resultNode.TransformArgsForCall(0)
				Expect(passedUserCtx).To(Equal(userCtx))
			})

			Context("if transform has no error", func() {
				transformedNode := &parsleyfakes.FakeNode{}
				transformedNode.TokenReturns("TRANSFORMED")

				BeforeEach(func() {
					resultNode.TransformReturns(transformedNode, nil)
				})
				It("returns the transformed node", func() {
					Expect(res).To(Equal(transformedNode))
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("if transform has an error", func() {
				transformErr := parsley.NewError(parsley.Pos(1), errors.New("transform error"))

				BeforeEach(func() {
					resultNode.TransformReturns(nil, transformErr)
				})
				It("returns the transformed node", func() {
					Expect(res).To(BeNil())
					Expect(err).To(MatchError("transform error at testpos"))
				})
			})

		})

		Context("if result is not transformable", func() {
			It("should return the original result", func() {
				Expect(res).To(Equal(parserRes))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("if transformation is disabled", func() {
		Context("if result can be transformed", func() {
			var resultNode *parsleyfakes.FakeTransformableNode

			BeforeEach(func() {
				resultNode = &parsleyfakes.FakeTransformableNode{}
				parserRes = resultNode
			})
			It("won't call the transform method", func() {
				Expect(resultNode.TransformCallCount()).To(BeZero())
			})
		})
	})

	Context("if static checking is enabled", func() {
		BeforeEach(func() {
			ctx.EnableStaticCheck()
		})

		Context("if result can be checked", func() {
			var resultNode *parsleyfakes.FakeStaticCheckableNode

			BeforeEach(func() {
				resultNode = &parsleyfakes.FakeStaticCheckableNode{}
				parserRes = resultNode
			})

			It("calls the StaticCheck function", func() {
				Expect(resultNode.StaticCheckCallCount()).To(Equal(1))
				passedUserCtx := resultNode.StaticCheckArgsForCall(0)
				Expect(passedUserCtx).To(Equal(userCtx))
			})

			Context("if static check has no error", func() {
				BeforeEach(func() {
					resultNode.StaticCheckReturns(nil)
				})
				It("returns no error", func() {
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("if static check has an error", func() {
				staticCheckErr := parsley.NewError(parsley.Pos(1), errors.New("static check error"))

				BeforeEach(func() {
					resultNode.StaticCheckReturns(staticCheckErr)
				})
				It("returns the transformed node", func() {
					Expect(err).To(MatchError("static check error at testpos"))
				})
			})

		})

		Context("if result is not static checkable", func() {
			It("should return no error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("if static checking is disabled", func() {
		Context("if result can be checked", func() {
			var resultNode *parsleyfakes.FakeStaticCheckableNode

			BeforeEach(func() {
				resultNode = &parsleyfakes.FakeStaticCheckableNode{}
				parserRes = resultNode
			})
			It("won't call the static check method", func() {
				Expect(resultNode.StaticCheckCallCount()).To(BeZero())
			})
		})
	})

	Context("if the parser has an error", func() {
		BeforeEach(func() {
			parserRes = nil
			err := &parsleyfakes.FakeError{}
			err.PosReturns(parsley.Pos(2))
			err.ErrorReturns("some error")
			parserErr = err
			ctx.SetError(parsley.NewError(parsley.Pos(1), errors.New("context error")))
		})
		It("should return an error", func() {
			Expect(val).To(BeNil())
			Expect(err).To(MatchError("failed to parse the input: some error at testpos"))
		})

		Context("if the context stores an error with a higher position", func() {
			BeforeEach(func() {
				ctx.SetError(parsley.NewError(parsley.Pos(3), errors.New("context error")))
			})
			It("should return an error", func() {
				Expect(val).To(BeNil())
				Expect(err).To(MatchError("failed to parse the input: context error at testpos"))
			})
		})
	})

})
