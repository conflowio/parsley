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
		err        error
		parserRes  *parsleyfakes.FakeNonTerminalNode
		parserErr  parsley.Error
		userCtx    interface{}
		nodeVal    interface{}
		nodeValErr parsley.Error
	)

	BeforeEach(func() {
		f := &parsleyfakes.FakeFile{}
		position := &parsleyfakes.FakePosition{}
		position.StringReturns("testpos")
		f.PositionReturns(position)

		fs := parsley.NewFileSet(f)

		r = &parsleyfakes.FakeReader{}
		r.PosReturns(parsley.Pos(1))
		p = &parsleyfakes.FakeParser{}
		parserRes = &parsleyfakes.FakeNonTerminalNode{}
		nodeVal = "value"
		userCtx = "context"
		parserErr = nil
		nodeValErr = nil

		ctx = parsley.NewContext(fs, r)
		ctx.SetUserContext(userCtx)
	})

	JustBeforeEach(func() {
		p.ParseReturns(parserRes, data.EmptyIntSet, parserErr)
		if parserRes != nil {
			parserRes.ValueReturns(nodeVal, nodeValErr)
		}
		val, err = parsley.Evaluate(ctx, p)
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

	It("gets the value from the node and passes the context", func() {
		Expect(parserRes.ValueCallCount()).To(Equal(1))
		Expect(parserRes.ValueArgsForCall(0)).To(Equal(userCtx))
	})

	It("should return the value of the node", func() {
		Expect(val).To(BeEquivalentTo(nodeVal))
		Expect(err).To(BeNil())
	})

	Context("if the parser has an error", func() {
		BeforeEach(func() {
			parserRes = nil
			err := &parsleyfakes.FakeError{}
			err.PosReturns(parsley.Pos(1))
			err.ErrorReturns("some error")
			parserErr = err
		})
		It("should return an error", func() {
			Expect(val).To(BeNil())
			Expect(err).To(MatchError("failed to parse the input: some error at testpos"))
		})
	})

	Context("if the node evaluation has an error", func() {
		BeforeEach(func() {
			nodeVal = nil
			err := &parsleyfakes.FakeError{}
			err.PosReturns(parsley.Pos(1))
			err.ErrorReturns("some error")
			nodeValErr = err
		})
		It("should return an error", func() {
			Expect(val).To(BeNil())
			Expect(err).To(MatchError("some error at testpos"))
		})
	})
})
