// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package combinator_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Let's define a parser which accepts integer or float numbers.
// The parser would return only the first match so in this case we have to put the float parser first.
func ExampleChoice() {
	p := combinator.Choice("number",
		terminal.Float(),
		terminal.Integer(),
	)
	r := text.NewReader(text.NewFile("example.file", []byte("1.23")))
	value, _ := parsley.Evaluate(parser.NewHistory(), r, combinator.Sentence(p), nil)
	fmt.Printf("%T %v\n", value, value)
	// Output: float64 1.23
}

var _ = Describe("Choice", func() {

	var (
		p                       *parser.NamedFunc
		h                       *parsleyfakes.FakeHistory
		r                       *parsleyfakes.FakeReader
		parsers                 []parsley.Parser
		p1, p2                  *parsleyfakes.FakeParser
		leftRecCtx              data.IntMap
		pos                     int
		cp, p1CP, p2CP          data.IntSet
		res, p1Res, p2Res       parsley.Node
		parserErr, p1Err, p2Err parsley.Error
		n1, n2                  *parsleyfakes.FakeNode
	)

	BeforeEach(func() {
		h = &parsleyfakes.FakeHistory{}
		r = &parsleyfakes.FakeReader{}
		p1 = &parsleyfakes.FakeParser{}
		p1.NameReturns("p1")
		p2 = &parsleyfakes.FakeParser{}
		p2.NameReturns("p2")
		leftRecCtx = data.EmptyIntMap
		parsers = []parsley.Parser{p1, p2}

		n1 = &parsleyfakes.FakeNode{}
		n1.TokenReturns("n1")
		n2 = &parsleyfakes.FakeNode{}
		n2.TokenReturns("n2")

		p1CP = data.EmptyIntSet
		p2CP = data.EmptyIntSet
		p1Res = nil
		p2Res = nil
		p1Err = nil
		p2Err = nil
		n1 = nil
		n2 = nil
	})

	JustBeforeEach(func() {
		p1.ParseReturnsOnCall(0, p1CP, p1Res, p1Err)
		p2.ParseReturnsOnCall(0, p2CP, p2Res, p2Err)

		p = combinator.Choice("test", parsers...)
		cp, res, parserErr = p.Parse(h, leftRecCtx, r, pos)
	})

	Context("when no parsers are given", func() {
		It("should panic", func() {
			Expect(func() { combinator.Choice("test") }).To(Panic())
		})
	})

	Context("when there is only one parser", func() {

		BeforeEach(func() {
			parsers = []parsley.Parser{p1}
			p1CP = data.NewIntSet(1)
			p1Res = n1
			p1Err = parsley.NewError(parsley.Pos(1), "some error")
		})

		It("should return the result of that parser", func() {
			Expect(cp).To(Equal(p1CP))
			Expect(res).To(Equal(p1Res))
			Expect(parserErr).To(Equal(p1Err))

			Expect(p1.ParseCallCount()).To(Equal(1))

			passedH, passedLeftRecCtx, passedR, passedPos := p1.ParseArgsForCall(0)
			Expect(passedH).To(BeEquivalentTo(h))
			Expect(passedLeftRecCtx).To(BeEquivalentTo(leftRecCtx))
			Expect(passedR).To(BeEquivalentTo(r))
			Expect(passedPos).To(Equal(pos))
		})
	})

	Context("with multiple parsers", func() {

		BeforeEach(func() {
			parsers = []parsley.Parser{p1, p2}
			p1CP = data.NewIntSet(1)
			p2CP = data.NewIntSet(2)
		})

		It("should merge the curtailing parsers", func() {
			Expect(cp).To(Equal(p1CP.Union(p2CP)))
		})

		Context("when there are multiple errors", func() {
			BeforeEach(func() {
				p1Err = parsley.NewError(parsley.Pos(1), "err1")
				p2Err = parsley.NewError(parsley.Pos(2), "err2")
			})
			It("should return with the error with the higher position", func() {
				Expect(parserErr).To(MatchError("err2"))
			})
		})

		Context("when no parsers match", func() {
			It("should return nil", func() {
				Expect(res).To(BeNil())
			})

			It("should call all parsers", func() {
				Expect(p1.ParseCallCount()).To(Equal(1))

				passedH, passedLeftRecCtx, passedR, passedPos := p1.ParseArgsForCall(0)
				Expect(passedH).To(BeEquivalentTo(h))
				Expect(passedLeftRecCtx).To(BeEquivalentTo(leftRecCtx))
				Expect(passedR).To(BeEquivalentTo(r))
				Expect(passedPos).To(Equal(pos))

				Expect(p2.ParseCallCount()).To(Equal(1))

				passedH, passedLeftRecCtx, passedR, passedPos = p2.ParseArgsForCall(0)
				Expect(passedH).To(BeEquivalentTo(h))
				Expect(passedLeftRecCtx).To(BeEquivalentTo(leftRecCtx))
				Expect(passedR).To(BeEquivalentTo(r))
				Expect(passedPos).To(Equal(pos))
			})
		})

		Context("when one parser matches", func() {
			BeforeEach(func() {
				p1Res = n1
			})
			It("should return the result of that parser", func() {
				Expect(res).To(Equal(p1Res))
			})
			It("should not call the remaining parsers", func() {
				Expect(p2.ParseCallCount()).To(Equal(0))
			})
		})
	})

})
