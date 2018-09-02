// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("Evaluate", func() {
	var (
		r          *parsleyfakes.FakeReader
		p          *parsleyfakes.FakeParser
		ctx        *parsley.Context
		val        interface{}
		err        parsley.Error
		parserRes  parsley.Node
		node       *parsleyfakes.FakeNode
		evalCtx    interface{}
		nodeVal    interface{}
		nodeValErr parsley.Error
	)

	BeforeEach(func() {
		r = &parsleyfakes.FakeReader{}
		ctx = parsley.NewContext(r)
		r.PosReturns(parsley.Pos(1))
		p = &parsleyfakes.FakeParser{}
		p.NameReturns("p1")
		node = &parsleyfakes.FakeNode{}
		parserRes = node
		nodeVal = "value"
		nodeValErr = nil
		evalCtx = "context"
	})

	JustBeforeEach(func() {
		node.ValueReturns(nodeVal, nodeValErr)
		p.ParseReturns(parserRes, data.EmptyIntSet)
		val, err = parsley.Evaluate(ctx, p, evalCtx)
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
	})

	It("gets the value from the node and passes the context", func() {
		Expect(node.ValueCallCount()).To(Equal(1))
		Expect(node.ValueArgsForCall(0)).To(Equal(evalCtx))
	})

	It("should return the value of the node", func() {
		Expect(val).To(BeEquivalentTo(nodeVal))
		Expect(err).To(BeNil())
	})

	Context("if the parser has an error", func() {
		BeforeEach(func() {
			parserRes = nil
			ctx.SetErrorf(parsley.Pos(1), "some error")
		})
		It("should return an error", func() {
			Expect(val).To(BeNil())
			Expect(err).To(MatchError("failed to parse the input: some error"))
			Expect(err.Pos()).To(Equal(parsley.Pos(1)))
		})
	})

	Context("if the node evaluation has an error", func() {
		BeforeEach(func() {
			nodeVal = nil
			nodeValErr = parsley.NewErrorf(parsley.Pos(1), "some error")
		})
		It("should return an error", func() {
			Expect(val).To(BeNil())
			Expect(err).To(MatchError("some error"))
			Expect(err.Pos()).To(Equal(parsley.Pos(1)))
		})
	})
})
