// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
		err, parserErr parsley.Error
		leftRectCtx    data.IntMap
		cp, parserCP   data.IntSet
	)

	JustBeforeEach(func() {
		fakep.ParseReturns(parserRes, parserCP, parserErr)
		f = text.NewFile("testfile", input)
		fs := parsley.NewFileSet(f)
		r = text.NewReader(f)
		ctx = parsley.NewContext(fs, r)
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
			input = []byte(" \t\n\fabc \t")
		})

		JustBeforeEach(func() {
			leftTrim := text.LeftTrim(p, wsMode)
			res, cp, err = leftTrim.Parse(ctx, leftRectCtx, pos)
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

		Context("when whitespace mode is forcing a new line", func() {
			BeforeEach(func() {
				wsMode = text.WsSpacesForceNl
			})

			It("should trim the spaces and new lines from the left", func() {
				_, _, passedPos := fakep.ParseArgsForCall(0)
				Expect(passedPos).To(Equal(pos + 4))
			})
		})

		Context("when whitespace mode is forcing a new line but no new line", func() {
			BeforeEach(func() {
				wsMode = text.WsSpacesForceNl
				pos = parsley.Pos(8)
			})

			It("should return an error", func() {
				Expect(err).To(MatchError(parsley.NewErrorf(parsley.Pos(10), "was expecting a new line")))
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
			res, cp, err = rightTrim.Parse(ctx, leftRectCtx, pos)
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
				p = terminal.Word("abc", "abc", "string")
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
					Expect(err).To(Equal(parserErr))
				})
			})

			Context("when whitespace mode is spaces", func() {
				BeforeEach(func() {
					wsMode = text.WsSpaces
				})

				It("should trim the spaces from the right", func() {
					Expect(err.Error()).To(Equal("some error"))
					Expect(err.Pos()).To(Equal(parsley.Pos(6)))
				})
			})

			Context("when whitespace mode is spaces and new lines", func() {
				BeforeEach(func() {
					wsMode = text.WsSpacesNl
				})

				It("should trim the spaces and new lines from the right", func() {
					Expect(err.Error()).To(Equal("some error"))
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
			res, cp, err = trim.Parse(ctx, leftRectCtx, pos)
		})

		Context("When there is result", func() {
			BeforeEach(func() {
				p = terminal.Word("abc", "abc", "string")
			})

			It("should trim the spaces and new lines from the right", func() {
				Expect(res).To(Equal(ast.NewTerminalNode("ABC", "abc", "string", parsley.Pos(5), parsley.Pos(12))))
			})
		})
	})
})
