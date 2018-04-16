package text_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Trim parsers", func() {

	var (
		r              *text.Reader
		f              *text.File
		input          []byte
		pos            parsley.Pos
		p              parsley.Parser
		fakep          *parsleyfakes.FakeParser
		wsMode         text.WsMode
		res, parserRes parsley.Node
		err, parserErr parsley.Error
		h              *parser.History
		leftRectCtx    data.IntMap
		cp, parserCP   data.IntSet
	)

	JustBeforeEach(func() {
		fakep.ParseReturns(parserRes, parserErr, parserCP)
		f = text.NewFile("testfile", input)
		r = text.NewReader(f)
	})

	BeforeEach(func() {
		fakep = &parsleyfakes.FakeParser{}
		p = fakep
		wsMode = text.WsNone
		pos = parsley.Pos(1)
		h = parser.NewHistory()
		leftRectCtx = data.EmptyIntMap.Inc(1)
		parserRes = &parsleyfakes.FakeNode{}
		parserErr = parsley.NewError(parsley.Pos(1), "some error")
		parserCP = data.EmptyIntSet.Insert(1)
	})

	var _ = Describe("LeftTrim", func() {
		BeforeEach(func() {
			input = []byte(" \t\n\fabc")
		})

		JustBeforeEach(func() {
			leftTrim := text.LeftTrim(p, wsMode)
			res, err, cp = leftTrim.Parse(h, leftRectCtx, r, pos)
		})

		It("calls the parser", func() {
			Expect(fakep.ParseCallCount()).To(Equal(1))
			passedHistory, passedLeftRecCtx, passedReader, _ := fakep.ParseArgsForCall(0)
			Expect(passedHistory).To(Equal(h))
			Expect(passedLeftRecCtx).To(Equal(leftRectCtx))
			Expect(passedReader).To(Equal(r))
		})

		It("returns the result of the parser", func() {
			Expect(res).To(BeIdenticalTo(parserRes))
			Expect(err).To(Equal(parserErr))
			Expect(cp).To(Equal(parserCP))
		})

		Context("when whitespace mode is no whitespaces", func() {
			BeforeEach(func() {
				wsMode = text.WsNone
			})

			It("skips no spaces", func() {
				_, _, _, passedPos := fakep.ParseArgsForCall(0)
				Expect(passedPos).To(Equal(pos))
			})
		})

		Context("when whitespace mode is spaces", func() {
			BeforeEach(func() {
				wsMode = text.WsSpaces
			})

			It("should trim the spaces from the left", func() {
				_, _, _, passedPos := fakep.ParseArgsForCall(0)
				Expect(passedPos).To(Equal(pos + 2))
			})
		})

		Context("when whitespace mode is spaces and new lines", func() {
			BeforeEach(func() {
				wsMode = text.WsSpacesNl
			})

			It("should trim the spaces and new lines from the left", func() {
				_, _, _, passedPos := fakep.ParseArgsForCall(0)
				Expect(passedPos).To(Equal(pos + 4))
			})
		})
	})

	var _ = Describe("RightTrim", func() {
		BeforeEach(func() {
			input = []byte("abc \t\n\f")
			parserRes = nil
			parserErr = nil
		})

		JustBeforeEach(func() {
			rightTrim := text.RightTrim(p, wsMode)
			res, err, cp = rightTrim.Parse(h, leftRectCtx, r, pos)
		})

		It("calls the parser", func() {
			Expect(fakep.ParseCallCount()).To(Equal(1))
			passedHistory, passedLeftRecCtx, passedReader, _ := fakep.ParseArgsForCall(0)
			Expect(passedHistory).To(Equal(h))
			Expect(passedLeftRecCtx).To(Equal(leftRectCtx))
			Expect(passedReader).To(Equal(r))
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
				parserErr = parsley.NewError(parsley.Pos(4), "some error")
			})

			Context("when whitespace mode is no whitespaces", func() {
				BeforeEach(func() {
					wsMode = text.WsNone
				})

				It("skips no spaces", func() {
					Expect(err).To(Equal(parserErr))
				})
			})

			Context("when whitespace mode is spaces", func() {
				BeforeEach(func() {
					wsMode = text.WsSpaces
				})

				It("should trim the spaces from the right", func() {
					Expect(err.Error()).To(Equal(parserErr.Error()))
					Expect(err.Pos()).To(Equal(parsley.Pos(6)))
				})
			})

			Context("when whitespace mode is spaces and new lines", func() {
				BeforeEach(func() {
					wsMode = text.WsSpacesNl
				})

				It("should trim the spaces and new lines from the right", func() {
					Expect(err.Error()).To(Equal(parserErr.Error()))
					Expect(err.Pos()).To(Equal(parsley.Pos(8)))
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
			res, err, cp = trim.Parse(h, leftRectCtx, r, pos)
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
