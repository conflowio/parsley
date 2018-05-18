package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("Nil", func() {

	var (
		p = parser.Nil()
		r *parsleyfakes.FakeReader
		h *parsleyfakes.FakeHistory
	)

	BeforeEach(func() {
		r = &parsleyfakes.FakeReader{}
		h = &parsleyfakes.FakeHistory{}
	})

	It("should have no name", func() {
		Expect(p.Name()).To(BeEmpty())
	})

	It("should return with an empty node", func() {
		res, err, curtailingParsers := p.Parse(h, data.EmptyIntMap, r, 1)
		Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
		Expect(res).To(Equal(ast.NilNode(1)))
		Expect(err).ToNot(HaveOccurred())
	})

})

var _ = Describe("End", func() {

	var (
		p = parser.End()
		r *parsleyfakes.FakeReader
		h *parsleyfakes.FakeHistory
	)

	BeforeEach(func() {
		r = &parsleyfakes.FakeReader{}
		h = &parsleyfakes.FakeHistory{}
	})

	It("should have a name", func() {
		Expect(p.Name()).ToNot(BeEmpty())
	})

	Context("when at the end of the input", func() {
		It("should return with an EOF node", func() {
			r.IsEOFReturns(true)
			r.PosReturns(parsley.Pos(2))
			res, err, curtailingParsers := p.Parse(h, data.EmptyIntMap, r, 2)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(res).To(Equal(ast.NewTerminalNode("EOF", nil, parsley.Pos(2), parsley.Pos(2))))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when not at the end of the input", func() {
		It("should return with a nil result", func() {
			r.IsEOFReturns(false)
			r.PosReturns(parsley.Pos(2))
			res, err, curtailingParsers := p.Parse(h, data.EmptyIntMap, r, 1)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(res).To(BeNil())
			Expect(err).To(BeNil())
		})
	})

})
