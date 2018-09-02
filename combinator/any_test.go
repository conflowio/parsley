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
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/combinator"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

// Let's define a parser which accepts integer or float numbers.
// The parser would return with all matches, so both 1 and 1.23.
func ExampleAny() {
	p := combinator.Any("number",
		terminal.Integer(),
		terminal.Float(),
	)
	r := text.NewReader(text.NewFile("example.file", []byte("1.23")))
	ctx := parsley.NewContext(r)
	value, _ := parsley.Evaluate(ctx, combinator.Sentence(p), nil)
	fmt.Printf("%T %v\n", value, value)

	// Output: float64 1.23
}

var _ = Describe("Any", func() {

	var (
		p                 *parser.NamedFunc
		r                 *parsleyfakes.FakeReader
		parsers           []parsley.Parser
		p1, p2            *parsleyfakes.FakeParser
		leftRecCtx        data.IntMap
		pos               parsley.Pos
		cp, p1CP, p2CP    data.IntSet
		res, p1Res, p2Res parsley.Node
		n1, n2            *parsleyfakes.FakeNode
		ctx               *parsley.Context
	)

	BeforeEach(func() {
		r = &parsleyfakes.FakeReader{}
		ctx = parsley.NewContext(r)
		p1 = &parsleyfakes.FakeParser{}
		p1.NameReturns("p1")
		p2 = &parsleyfakes.FakeParser{}
		p2.NameReturns("p2")
		leftRecCtx = data.EmptyIntMap
		parsers = []parsley.Parser{p1, p2}
		pos = parsley.Pos(1)

		n1 = &parsleyfakes.FakeNode{}
		n1.TokenReturns("n1")
		n2 = &parsleyfakes.FakeNode{}
		n2.TokenReturns("n2")

		p1CP = data.EmptyIntSet
		p2CP = data.EmptyIntSet
		p1Res = nil
		p2Res = nil
		n1 = nil
		n2 = nil
	})

	JustBeforeEach(func() {
		p1.ParseReturnsOnCall(0, p1Res, p1CP)
		p2.ParseReturnsOnCall(0, p2Res, p2CP)

		p = combinator.Any("test", parsers...)
		res, cp = p.Parse(ctx, leftRecCtx, pos)
	})

	Context("when no parsers are given", func() {
		It("should panic", func() {
			Expect(func() { combinator.Any("test") }).To(Panic())
		})
	})

	Context("when there is only one parser", func() {

		BeforeEach(func() {
			parsers = []parsley.Parser{p1}
			p1CP = data.NewIntSet(1)
			p1Res = n1
		})

		It("should return the result of that parser", func() {
			Expect(cp).To(Equal(p1CP))
			Expect(res).To(Equal(p1Res))

			Expect(p1.ParseCallCount()).To(Equal(1))

			passedCtx, passedLeftRecCtx, passedPos := p1.ParseArgsForCall(0)
			Expect(passedCtx).To(BeEquivalentTo(ctx))
			Expect(passedLeftRecCtx).To(BeEquivalentTo(leftRecCtx))
			Expect(passedPos).To(Equal(pos))
		})
	})

	Context("with multiple parsers", func() {

		BeforeEach(func() {
			parsers = []parsley.Parser{p1, p2}
			p1CP = data.NewIntSet(1)
			p2CP = data.NewIntSet(2)
		})

		It("should call all parsers", func() {
			Expect(p1.ParseCallCount()).To(Equal(1))

			passedCtx, passedLeftRecCtx, passedPos := p1.ParseArgsForCall(0)
			Expect(passedCtx).To(BeEquivalentTo(ctx))
			Expect(passedLeftRecCtx).To(BeEquivalentTo(leftRecCtx))
			Expect(passedPos).To(Equal(pos))

			Expect(p2.ParseCallCount()).To(Equal(1))

			passedCtx, passedLeftRecCtx, passedPos = p2.ParseArgsForCall(0)
			Expect(passedCtx).To(BeEquivalentTo(ctx))
			Expect(passedLeftRecCtx).To(BeEquivalentTo(leftRecCtx))
			Expect(passedPos).To(Equal(pos))
		})

		It("should merge the curtailing parsers", func() {
			Expect(cp).To(Equal(p1CP.Union(p2CP)))
		})

		Context("when no parsers match", func() {
			It("should return nil", func() {
				Expect(res).To(BeNil())
			})
		})

		Context("when one parser matches", func() {
			BeforeEach(func() {
				p1Res = n1
			})
			It("should return the result of that parser", func() {
				Expect(res).To(Equal(p1Res))
			})
		})

		Context("when multiple parsers match", func() {
			BeforeEach(func() {
				p1Res = n1
				p2Res = n2
			})
			It("should return with all results", func() {
				Expect(res).To(Equal(ast.NodeList([]parsley.Node{n1, n2})))
			})
		})
	})

})
