// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley_test

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

// Let's define a simple parser which is able to parse adding two integers.
func ExampleParse() {
	sum := ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		value0, _ := nodes[0].Value(ctx)
		value1, _ := nodes[2].Value(ctx)
		return value0.(int) + value1.(int), nil
	})

	p := combinator.Seq("ADD", "addition",
		terminal.Integer(),
		terminal.Rune('+'),
		terminal.Integer(),
	).Bind(sum)

	r := text.NewReader(text.NewFile("example.file", []byte("1+2")))
	node, err := parsley.Parse(parser.NewHistory(), r, combinator.Sentence(p))
	if err != nil {
		panic(err)
	}

	value, err2 := node.Value(nil)
	if err2 != nil {
		panic(err2)
	}
	fmt.Printf("Result: %d\n", value.(int))
	// Output: Result: 3
}

// Let's define a simple parser which is able to parse adding two integers.
func ExampleEvaluate() {
	sum := ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
		value0, _ := nodes[0].Value(ctx)
		value1, _ := nodes[2].Value(ctx)
		return value0.(int) + value1.(int), nil
	})

	p := combinator.Seq("ADD", "addition",
		terminal.Integer(),
		terminal.Rune('+'),
		terminal.Integer(),
	).Bind(sum)

	r := text.NewReader(text.NewFile("example.file", []byte("1+2")))
	value, err := parsley.Evaluate(parser.NewHistory(), r, combinator.Sentence(p), nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %d\n", value.(int))
	// Output: Result: 3
}

var _ = Describe("Parse", func() {
	var (
		h         *parsleyfakes.FakeHistory
		r         *parsleyfakes.FakeReader
		p         *parsleyfakes.FakeParser
		res       parsley.Node
		err       parsley.Error
		parserRes parsley.Node
		parserErr parsley.Error
	)

	BeforeEach(func() {
		h = &parsleyfakes.FakeHistory{}
		r = &parsleyfakes.FakeReader{}
		r.PosReturns(parsley.Pos(1))
		p = &parsleyfakes.FakeParser{}
		p.NameReturns("p1")
		parserRes = &parsleyfakes.FakeNode{}
		parserErr = nil
	})

	JustBeforeEach(func() {
		p.ParseReturns(parserRes, parserErr, data.EmptyIntSet)
		res, err = parsley.Parse(h, r, p)
	})

	It("gets the zero position from the reader", func() {
		Expect(r.PosCallCount()).To(Equal(1))
		Expect(r.PosArgsForCall(0)).To(Equal(0))
	})

	It("calls the parser", func() {
		Expect(p.ParseCallCount()).To(Equal(1))
		passedHistory, passedLeftRecCtx, passedReader, passedPos := p.ParseArgsForCall(0)
		Expect(passedHistory).To(BeEquivalentTo(h))
		Expect(passedLeftRecCtx).To(BeEquivalentTo(data.EmptyIntMap))
		Expect(passedReader).To(BeEquivalentTo(r))
		Expect(passedPos).To(Equal(parsley.Pos(1)))
	})

	It("should return the result of the parser", func() {
		Expect(res).To(BeEquivalentTo(parserRes))
		Expect(err).To(BeNil())
	})

	Context("when parser returned with no result and an error", func() {
		BeforeEach(func() {
			parserRes = nil
			parserErr = parsley.NewErrorf(parsley.Pos(1), "some error")
		})
		It("should return the error", func() {
			Expect(res).To(BeNil())
			Expect(err).To(MatchError("failed to parse the input: some error"))
			Expect(err.Pos()).To(Equal(parsley.Pos(1)))
		})
	})

	Context("when parser returned with no result and no error", func() {
		BeforeEach(func() {
			parserRes = nil
			parserErr = nil
		})
		It("should return with an error saying expecting the parser's name", func() {
			Expect(res).To(BeNil())
			Expect(err).To(MatchError("failed to parse the input: was expecting p1"))
			Expect(err.Pos()).To(Equal(parsley.Pos(1)))
		})
	})

	Context("when the result node is nil", func() {
		BeforeEach(func() {
			parserRes = nil
			parserErr = nil
		})
		It("should return with an error saying expecting the parser's name", func() {
			Expect(res).To(BeNil())
			Expect(err).To(MatchError("failed to parse the input: was expecting p1"))
			Expect(err.Pos()).To(Equal(parsley.Pos(1)))
		})
	})
})

var _ = Describe("Evaluate", func() {
	var (
		h          *parsleyfakes.FakeHistory
		r          *parsleyfakes.FakeReader
		p          *parsleyfakes.FakeParser
		val        interface{}
		err        parsley.Error
		parserRes  parsley.Node
		node       *parsleyfakes.FakeNode
		parserErr  parsley.Error
		ctx        interface{}
		nodeVal    interface{}
		nodeValErr parsley.Error
	)

	BeforeEach(func() {
		h = &parsleyfakes.FakeHistory{}
		r = &parsleyfakes.FakeReader{}
		r.PosReturns(parsley.Pos(1))
		p = &parsleyfakes.FakeParser{}
		p.NameReturns("p1")
		node = &parsleyfakes.FakeNode{}
		parserRes = node
		parserErr = nil
		nodeVal = "value"
		nodeValErr = nil
		ctx = "context"
	})

	JustBeforeEach(func() {
		node.ValueReturns(nodeVal, nodeValErr)
		p.ParseReturns(parserRes, parserErr, data.EmptyIntSet)
		val, err = parsley.Evaluate(h, r, p, ctx)
	})

	It("gets the zero position from the reader", func() {
		Expect(r.PosCallCount()).To(Equal(1))
		Expect(r.PosArgsForCall(0)).To(Equal(0))
	})

	It("calls the parser", func() {
		Expect(p.ParseCallCount()).To(Equal(1))
		passedHistory, passedLeftRecCtx, passedReader, passedPos := p.ParseArgsForCall(0)
		Expect(passedHistory).To(BeEquivalentTo(h))
		Expect(passedLeftRecCtx).To(BeEquivalentTo(data.EmptyIntMap))
		Expect(passedReader).To(BeEquivalentTo(r))
		Expect(passedPos).To(Equal(parsley.Pos(1)))
	})

	It("gets the value from the node and passes the context", func() {
		Expect(node.ValueCallCount()).To(Equal(1))
		Expect(node.ValueArgsForCall(0)).To(Equal(ctx))
	})

	It("should return the value of the node", func() {
		Expect(val).To(BeEquivalentTo(nodeVal))
		Expect(err).To(BeNil())
	})

	Context("if the parser has an error", func() {
		BeforeEach(func() {
			parserRes = nil
			parserErr = parsley.NewErrorf(parsley.Pos(1), "some error")
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

var _ = Describe("Parsley", func() {

	It("should handle direct left recursion", func() {
		input := "abbbbbbbbbbbbbbbbbbb"

		concat := ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
			s := ""
			for _, node := range nodes {
				val, err := node.Value(ctx)
				if err != nil {
					return nil, err
				}
				if runeVal, ok := val.(rune); ok {
					s = s + string(runeVal)
				} else {
					s = s + val.(string)
				}
			}
			return s, nil
		})

		var p parser.NamedFunc
		p = *combinator.Memoize(combinator.Any("a or ab",
			combinator.Seq("AB", "a or ab",
				&p,
				terminal.Rune('b'),
			).Bind(concat),
			terminal.Rune('a'),
		))

		h := parser.NewHistory()
		f := text.NewFile("testfile", []byte(input))
		r := text.NewReader(f)
		result, err := parsley.Evaluate(h, r, combinator.Sentence(&p), nil)

		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(input))
		Expect(h.CallCount()).To(Equal(299))
	})

	It("should handle highly ambiguous left-recursive grammar", func() {
		input := "1+2+3+4+5+6+7+8+9+10"

		add := ast.InterpreterFunc(func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
			value0, _ := nodes[0].Value(ctx)
			value1, _ := nodes[2].Value(ctx)
			return value0.(int) + value1.(int), nil
		})

		var p parser.NamedFunc
		value := combinator.Memoize(combinator.Any("value",
			terminal.Integer(),
			&p,
		))

		p = *combinator.Memoize(combinator.Seq("ADD", "addition",
			value,
			terminal.Rune('+'),
			value,
		).Bind(add))

		h := parser.NewHistory()
		f := text.NewFile("testfile", []byte(input))
		r := text.NewReader(f)
		result, err := parsley.Evaluate(h, r, combinator.Sentence(&p), nil)

		Expect(err).ToNot(HaveOccurred())
		Expect(result).To(Equal(55))
		Expect(h.CallCount()).To(Equal(237770))

	})
})
