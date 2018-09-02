package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("History", func() {
	var (
		h *parser.History
	)

	BeforeEach(func() {
		h = parser.NewHistory()
	})

	Describe("SaveResult/GetResult", func() {
		var (
			parserIndex           int
			pos                   parsley.Pos
			res, res1, res2, res3 *parsley.Result
			found                 bool
			leftRecCtx            data.IntMap
		)

		BeforeEach(func() {
			res1 = &parsley.Result{
				CurtailingParsers: data.EmptyIntSet,
				LeftRecCtx:        data.EmptyIntMap,
				Node:              &parsleyfakes.FakeNode{},
			}
			res2 = &parsley.Result{
				CurtailingParsers: data.EmptyIntSet,
				LeftRecCtx:        data.EmptyIntMap,
				Node:              &parsleyfakes.FakeNode{},
			}
			res3 = &parsley.Result{
				CurtailingParsers: data.EmptyIntSet,
				LeftRecCtx:        data.EmptyIntMap,
				Node:              &parsleyfakes.FakeNode{},
			}

			leftRecCtx = data.EmptyIntMap
			parserIndex = 1
			pos = parsley.Pos(1)
		})

		JustBeforeEach(func() {
			h.SaveResult(1, parsley.Pos(1), res1)
			h.SaveResult(1, parsley.Pos(2), res2)
			h.SaveResult(2, parsley.Pos(1), res3)
			res, found = h.GetResult(parserIndex, pos, leftRecCtx)
		})

		It("should return previously saved result", func() {
			Expect(res).To(BeEquivalentTo(res1))
			Expect(found).To(BeTrue())
		})

		Context("when getting a result for the same parser but different position", func() {
			BeforeEach(func() {
				pos = parsley.Pos(2)
			})
			It("should return a different result", func() {
				Expect(res).To(BeEquivalentTo(res2))
				Expect(found).To(BeTrue())
			})
		})

		Context("when getting a result for a different parser but the same position", func() {
			BeforeEach(func() {
				parserIndex = 2
			})
			It("should return a different result", func() {
				Expect(res).To(BeEquivalentTo(res3))
				Expect(found).To(BeTrue())
			})
		})

		Context("when there is no result for the parser", func() {
			BeforeEach(func() {
				parserIndex = 99
			})
			It("should return nil", func() {
				Expect(res).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when there is no result for the position", func() {
			BeforeEach(func() {
				pos = parsley.Pos(99)
			})
			It("should return nil", func() {
				Expect(res).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when there is no result for the position", func() {
			BeforeEach(func() {
				pos = parsley.Pos(99)
			})
			It("should return nil", func() {
				Expect(res).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when in the saved result any of the left-rec counts are higher for a parser", func() {
			BeforeEach(func() {
				res1.LeftRecCtx = data.NewIntMap(map[int]int{1: 2, 2: 1})
				leftRecCtx = data.NewIntMap(map[int]int{1: 1, 2: 1})
			})
			It("should not be returned", func() {
				Expect(res).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})
	})

	Describe("RegisterCall/CallCount", func() {
		It("should register a call count", func() {
			Expect(h.CallCount()).To(Equal(0))
			h.RegisterCall()
			Expect(h.CallCount()).To(Equal(1))
		})
	})

})
