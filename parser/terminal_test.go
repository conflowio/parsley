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

var _ = Describe("Empty", func() {

	var (
		p = parser.Empty()
		r *parsleyfakes.FakeReader
		h *parsleyfakes.FakeHistory
	)

	BeforeEach(func() {
		r = &parsleyfakes.FakeReader{}
		h = &parsleyfakes.FakeHistory{}
	})

	It("should return with a nil node", func() {
		curtailingParsers, res, err := p.Parse(h, data.EmptyIntMap, r, 0)
		Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
		Expect(res).To(Equal([]parsley.Node{nil}))
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

	Context("when at the end of the input", func() {
		It("should return with an EOF node", func() {
			r.IsEOFReturns(true)
			r.PosReturns(parsley.Pos(2))
			curtailingParsers, res, err := p.Parse(h, data.EmptyIntMap, r, 1)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(res).To(Equal([]parsley.Node{ast.NewTerminalNode("EOF", nil, parsley.Pos(2), 1)}))
			Expect(err).ToNot(HaveOccurred())

			Expect(r.PosCallCount()).To(Equal(1))
			pos := r.PosArgsForCall(0)
			Expect(pos).To(Equal(1))
		})
	})

	Context("when not at the end of the input", func() {
		It("should return with an error", func() {
			r.IsEOFReturns(false)
			r.PosReturns(parsley.Pos(2))
			curtailingParsers, res, err := p.Parse(h, data.EmptyIntMap, r, 1)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(res).To(BeNil())
			Expect(err).To(MatchError("was expecting the end of input"))

			Expect(r.PosCallCount()).To(Equal(1))
			pos := r.PosArgsForCall(0)
			Expect(pos).To(Equal(1))
		})
	})

})
