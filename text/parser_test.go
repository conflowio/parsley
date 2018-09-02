package text_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Trim parsers", func() {

	var (
		r              *text.Reader
		f              *text.File
		ctx            *parsley.Context
		input          []byte
		pos            parsley.Pos
		p              parsley.Parser
		fakep          *parsleyfakes.FakeParser
		wsMode         text.WsMode
		res, parserRes parsley.Node
		parserErr      parsley.Error
		leftRectCtx    data.IntMap
		cp, parserCP   data.IntSet
	)

	JustBeforeEach(func() {
		fakep.ParseReturns(parserRes, parserCP)
		f = text.NewFile("testfile", input)
		r = text.NewReader(f)
		ctx = parsley.NewContext(r, nil)
		ctx.OverrideError(parserErr)
	})

	BeforeEach(func() {
		fakep = &parsleyfakes.FakeParser{}
		p = fakep
		wsMode = text.WsNone
		pos = parsley.Pos(1)
		leftRectCtx = data.EmptyIntMap.Inc(1)
		parserRes = &parsleyfakes.FakeNode{}
		parserErr = parsley.NewErrorf(parsley.Pos(1), "some error")
		parserCP = data.EmptyIntSet.Insert(1)
	})

	var _ = Describe("LeftTrim", func() {
		BeforeEach(func() {
			input = []byte(" \t\n\fabc")
		})

		JustBeforeEach(func() {
			leftTrim := text.LeftTrim(p, wsMode)
			res, cp = leftTrim.Parse(ctx, leftRectCtx, pos)
		})

		It("calls the parser", func() {
			Expect(fakep.ParseCallCount()).To(Equal(1))
			passedCtx, passedLeftRecCtx, _ := fakep.ParseArgsForCall(0)
			Expect(passedCtx).To(Equal(ctx))
			Expect(passedLeftRecCtx).To(Equal(leftRectCtx))
		})

		It("returns the result of the parser", func() {
			Expect(res).To(BeIdenticalTo(parserRes))
			Expect(cp).To(Equal(parserCP))
		})

		Context("when whitespace mode is no whitespaces", func() {
			BeforeEach(func() {
				wsMode = text.WsNone
			})

			It("skips no spaces", func() {
				_, _, passedPos := fakep.ParseArgsForCall(0)
				Expect(passedPos).To(Equal(pos))
			})
		})

		Context("when whitespace mode is spaces", func() {
			BeforeEach(func() {
				wsMode = text.WsSpaces
			})

			It("should trim the spaces from the left", func() {
				_, _, passedPos := fakep.ParseArgsForCall(0)
				Expect(passedPos).To(Equal(pos + 2))
			})
		})

		Context("when whitespace mode is spaces and new lines", func() {
			BeforeEach(func() {
				wsMode = text.WsSpacesNl
			})

			It("should trim the spaces and new lines from the left", func() {
				_, _, passedPos := fakep.ParseArgsForCall(0)
				Expect(passedPos).To(Equal(pos + 4))
			})
		})
	})

	var _ = Describe("RightTrim", func() {
		BeforeEach(func() {
			input = []byte("abc \t\n\f")
			parserRes = nil
		})

		JustBeforeEach(func() {
			rightTrim := text.RightTrim(p, wsMode)
			res, cp = rightTrim.Parse(ctx, leftRectCtx, pos)
		})

		It("calls the parser", func() {
			Expect(fakep.ParseCallCount()).To(Equal(1))
			passedCtx, passedLeftRecCtx, _ := fakep.ParseArgsForCall(0)
			Expect(passedCtx).To(Equal(ctx))
			Expect(passedLeftRecCtx).To(Equal(leftRectCtx))
		})

		It("returns with the curtailing parsers", func() {
			Expect(cp).To(Equal(parserCP))
		})

		Context("When there is result", func() {
			BeforeEach(func() {
				p = terminal.Word("abc", "abc")
			})

			Context("when whitespace mode is no whitespaces", func() {
				BeforeEach(func() {
					wsMode = text.WsNone
				})

				It("skips no spaces", func() {
					Expect(res.ReaderPos()).To(Equal(parsley.Pos(4)))
				})
			})

			Context("when whitespace mode is spaces", func() {
				BeforeEach(func() {
					wsMode = text.WsSpaces
				})

				It("should trim the spaces from the right", func() {
					Expect(res.ReaderPos()).To(Equal(parsley.Pos(6)))
				})
			})

			Context("when whitespace mode is spaces and new lines", func() {
				BeforeEach(func() {
					wsMode = text.WsSpacesNl
				})

				It("should trim the spaces and new lines from the right", func() {
					Expect(res.ReaderPos()).To(Equal(parsley.Pos(8)))
				})
			})
		})

		Context("When there is error", func() {
			BeforeEach(func() {
				parserRes = nil
				parserErr = parsley.NewErrorf(parsley.Pos(4), "some error")
			})

			Context("when whitespace mode is no whitespaces", func() {
				BeforeEach(func() {
					wsMode = text.WsNone
				})

				It("skips no spaces", func() {
					Expect(ctx.Error()).To(Equal(parserErr))
				})
			})

			Context("when whitespace mode is spaces", func() {
				BeforeEach(func() {
					wsMode = text.WsSpaces
				})

				It("should trim the spaces from the right", func() {
					Expect(ctx.Error().Error()).To(Equal("some error"))
					Expect(ctx.Error().Pos()).To(Equal(parsley.Pos(6)))
				})
			})

			Context("when whitespace mode is spaces and new lines", func() {
				BeforeEach(func() {
					wsMode = text.WsSpacesNl
				})

				It("should trim the spaces and new lines from the right", func() {
					Expect(ctx.Error().Error()).To(Equal("some error"))
					Expect(ctx.Error().Pos()).To(Equal(parsley.Pos(8)))
				})
			})
		})
	})

	var _ = Describe("Trim", func() {
		BeforeEach(func() {
			input = []byte(" \t\n\fabc \t\n\f")
		})

		JustBeforeEach(func() {
			trim := text.Trim(p)
			res, cp = trim.Parse(ctx, leftRectCtx, pos)
		})

		Context("When there is result", func() {
			BeforeEach(func() {
				p = terminal.Word("abc", "abc")
			})

			It("should trim the spaces and new lines from the right", func() {
				Expect(res).To(Equal(ast.NewTerminalNode("WORD", "abc", parsley.Pos(5), parsley.Pos(12))))
			})
		})
	})
})
