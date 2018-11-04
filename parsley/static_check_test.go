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

var _ = Describe("StaticCheck", func() {
	var (
		r         *parsleyfakes.FakeReader
		p         *parsleyfakes.FakeParser
		ctx       *parsley.Context
		parserRes *parsleyfakes.FakeStaticCheckableNode
		parserErr parsley.Error
		evalCtx   interface{}
		checkErr  parsley.Error
		err       error
	)

	BeforeEach(func() {
		f := &parsleyfakes.FakeFile{}
		position := &parsleyfakes.FakePosition{}
		position.StringReturns("testpos")
		f.PositionReturns(position)

		fs := parsley.NewFileSet(f)

		r = &parsleyfakes.FakeReader{}
		ctx = parsley.NewContext(fs, r)
		r.PosReturns(parsley.Pos(1))
		p = &parsleyfakes.FakeParser{}
		parserRes = &parsleyfakes.FakeStaticCheckableNode{}
		evalCtx = "context"
		parserErr = nil
		checkErr = nil
		err = nil
	})

	JustBeforeEach(func() {
		p.ParseReturns(parserRes, data.EmptyIntSet, parserErr)
		if parserRes != nil {
			parserRes.StaticCheckReturns(checkErr)
		}
		err = parsley.StaticCheck(ctx, p, evalCtx)
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

	Context("if the parser has an error", func() {
		BeforeEach(func() {
			parserRes = nil
			err := &parsleyfakes.FakeError{}
			err.PosReturns(parsley.Pos(1))
			err.ErrorReturns("some error")
			parserErr = err
		})
		It("should return an error", func() {
			Expect(err).To(MatchError("failed to parse the input: some error at testpos"))
		})
	})

	Context("if the static check has an error", func() {
		BeforeEach(func() {
			err := &parsleyfakes.FakeError{}
			err.PosReturns(parsley.Pos(1))
			err.ErrorReturns("some error")
			checkErr = err
		})
		It("should return an error", func() {
			Expect(err).To(MatchError("some error at testpos"))
		})
	})
})
