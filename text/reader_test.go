// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package text_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
)

var _ = Describe("Reader", func() {

	const (
		input         = "abc def"
		inputWithUTF8 = "🍕 and 🍺"
	)

	var (
		r    *text.Reader
		f    *text.File
		data []byte
	)

	BeforeEach(func() {
		data = []byte(input)
	})

	JustBeforeEach(func() {
		f = text.NewFile("testfile", data)
		r = text.NewReader(f)
	})

	Describe("ReadRune()", func() {
		It("should match the next ASCII rune", func() {
			pos, found := r.ReadRune(f.Pos(0), 'a')
			Expect(pos).To(Equal(f.Pos(1)))
			Expect(found).To(BeTrue())
		})

		It("should not match a different rune", func() {
			pos, found := r.ReadRune(f.Pos(0), 'b')
			Expect(pos).To(Equal(f.Pos(0)))
			Expect(found).To(BeFalse())
		})

		It("should match an ASCII rune at the end", func() {
			pos, found := r.ReadRune(f.Pos(6), 'f')
			Expect(pos).To(Equal(f.Pos(7)))
			Expect(found).To(BeTrue())
		})

		Context("When input contains UTF8", func() {
			BeforeEach(func() {
				data = []byte(inputWithUTF8)
			})

			It("should match the next UTF8 rune", func() {
				pos, found := r.ReadRune(f.Pos(0), '🍕')
				Expect(pos).To(Equal(f.Pos(4)))
				Expect(found).To(BeTrue())
			})

			It("should match a UTF8 rune at the end", func() {
				pos, found := r.ReadRune(f.Pos(9), '🍺')
				Expect(pos).To(Equal(f.Pos(13)))
				Expect(found).To(BeTrue())
			})
		})

		Context("at the end of the file", func() {
			It("should not match anything", func() {
				pos, found := r.ReadRune(f.Pos(7), 'f')
				Expect(pos).To(Equal(f.Pos(7)))
				Expect(found).To(BeFalse())
			})
		})
	})

	Describe("MatchString()", func() {
		Context("when called with empty string", func() {
			It("should panic", func() {
				Expect(func() { r.MatchString(f.Pos(0), "") }).To(Panic())
			})
		})

		It("should match the a substring", func() {
			pos, found := r.MatchString(f.Pos(0), "ab")
			Expect(pos).To(Equal(f.Pos(2)))
			Expect(found).To(BeTrue())
		})

		It("should not match a different substring", func() {
			pos, found := r.MatchString(f.Pos(0), "ac")
			Expect(pos).To(Equal(f.Pos(0)))
			Expect(found).To(BeFalse())
		})

		It("should match a substring at the end", func() {
			pos, found := r.MatchString(f.Pos(5), "ef")
			Expect(pos).To(Equal(f.Pos(7)))
			Expect(found).To(BeTrue())
		})

		It("should only match the full substring at the end", func() {
			pos, found := r.MatchString(f.Pos(5), "efg")
			Expect(pos).To(Equal(f.Pos(5)))
			Expect(found).To(BeFalse())
		})

		Context("at the end of the file", func() {
			It("should not match anything", func() {
				pos, found := r.MatchString(f.Pos(7), "x")
				Expect(pos).To(Equal(f.Pos(7)))
				Expect(found).To(BeFalse())
			})
		})

		Context("When input contains UTF8", func() {
			BeforeEach(func() {
				data = []byte(inputWithUTF8)
			})

			It("should match the a substring", func() {
				pos, found := r.MatchString(f.Pos(0), "🍕 and")
				Expect(pos).To(Equal(f.Pos(8)))
				Expect(found).To(BeTrue())
			})

			It("should not match a different substring", func() {
				pos, found := r.MatchString(f.Pos(0), "🍕 not")
				Expect(pos).To(Equal(f.Pos(0)))
				Expect(found).To(BeFalse())
			})

			It("should match a substring at the end", func() {
				pos, found := r.MatchString(f.Pos(5), "and 🍺")
				Expect(pos).To(Equal(f.Pos(13)))
				Expect(found).To(BeTrue())
			})

			It("should only match the full substring at the end", func() {
				pos, found := r.MatchString(f.Pos(5), "and 🍺 s")
				Expect(pos).To(Equal(f.Pos(5)))
				Expect(found).To(BeFalse())
			})
		})
	})

	Describe("MatchWord()", func() {
		Context("when called with empty string", func() {
			It("should panic", func() {
				Expect(func() { r.MatchWord(f.Pos(0), "") }).To(Panic())
			})
		})

		It("should match the full word", func() {
			pos, found := r.MatchWord(f.Pos(0), "abc")
			Expect(pos).To(Equal(f.Pos(3)))
			Expect(found).To(BeTrue())
		})

		It("should not match a partial word", func() {
			pos, found := r.MatchWord(f.Pos(0), "ab")
			Expect(pos).To(Equal(f.Pos(0)))
			Expect(found).To(BeFalse())
		})

		It("should not match the different word", func() {
			pos, found := r.MatchWord(f.Pos(0), "abd")
			Expect(pos).To(Equal(f.Pos(0)))
			Expect(found).To(BeFalse())
		})

		It("should match a word at the end", func() {
			pos, found := r.MatchWord(f.Pos(4), "def")
			Expect(pos).To(Equal(f.Pos(7)))
			Expect(found).To(BeTrue())
		})

		It("should only match the full word at the end", func() {
			pos, found := r.MatchWord(f.Pos(4), "defg")
			Expect(pos).To(Equal(f.Pos(4)))
			Expect(found).To(BeFalse())
		})

		Context("at the end of the file", func() {
			It("should not match anything", func() {
				pos, found := r.MatchWord(f.Pos(7), "x")
				Expect(pos).To(Equal(f.Pos(7)))
				Expect(found).To(BeFalse())
			})
		})

		Context("When input contains UTF8", func() {
			BeforeEach(func() {
				data = []byte(inputWithUTF8)
			})

			It("should panic", func() {
				Expect(func() { r.MatchWord(f.Pos(0), "🍕 and") }).To(Panic())
			})
		})
	})

	Describe("ReadRegexp()", func() {
		Context("when matches an empty string", func() {
			It("should panic", func() {
				Expect(func() { r.ReadRegexp(f.Pos(0), "x?") }).To(Panic())
			})
		})

		It("should match the regexp", func() {
			pos, match := r.ReadRegexp(f.Pos(0), "a+b+x?")
			Expect(pos).To(Equal(f.Pos(2)))
			Expect(match).To(Equal([]byte("ab")))
		})

		It("should not match a non-matching regexp", func() {
			pos, match := r.ReadRegexp(f.Pos(0), "ac+")
			Expect(pos).To(Equal(f.Pos(0)))
			Expect(match).To(BeNil())
		})

		It("should match a regexp at the end", func() {
			pos, match := r.ReadRegexp(f.Pos(5), "ef+")
			Expect(pos).To(Equal(f.Pos(7)))
			Expect(match).To(Equal([]byte("ef")))
		})

		It("should only match the full match at the end", func() {
			pos, match := r.ReadRegexp(f.Pos(5), "efg+")
			Expect(pos).To(Equal(f.Pos(5)))
			Expect(match).To(BeNil())
		})

		Context("at the end of the file", func() {
			It("should not match anything", func() {
				pos, match := r.ReadRegexp(f.Pos(7), "x+")
				Expect(pos).To(Equal(f.Pos(7)))
				Expect(match).To(BeNil())
			})
		})

		Context("When input contains UTF8", func() {
			BeforeEach(func() {
				data = []byte(inputWithUTF8)
			})

			It("should match the regexp", func() {
				pos, match := r.ReadRegexp(f.Pos(0), ".* and")
				Expect(pos).To(Equal(f.Pos(8)))
				Expect(match).To(Equal([]byte("🍕 and")))
			})

			It("should not match a non-matching regexp", func() {
				pos, match := r.ReadRegexp(f.Pos(0), ".* not")
				Expect(pos).To(Equal(f.Pos(0)))
				Expect(match).To(BeNil())
			})

			It("should match a regexp at the end", func() {
				pos, match := r.ReadRegexp(f.Pos(5), "and .*")
				Expect(pos).To(Equal(f.Pos(13)))
				Expect(match).To(Equal([]byte("and 🍺")))
			})

			It("should only match the full match at the end", func() {
				pos, match := r.ReadRegexp(f.Pos(5), "and .*s")
				Expect(pos).To(Equal(f.Pos(5)))
				Expect(match).To(BeNil())
			})
		})
	})

	Describe("ReadRegexpSubmatch()", func() {
		Context("when matches an empty string", func() {
			It("should panic", func() {
				Expect(func() { r.ReadRegexpSubmatch(f.Pos(0), "x?") }).To(Panic())
			})
		})

		It("should match the regexp", func() {
			pos, match := r.ReadRegexpSubmatch(f.Pos(0), "(a+)b+x?")
			Expect(pos).To(Equal(f.Pos(2)))
			Expect(match).To(Equal([][]byte{
				[]byte("ab"),
				[]byte("a"),
			}))
		})

		It("should not match a non-matching regexp", func() {
			pos, match := r.ReadRegexpSubmatch(f.Pos(0), "(a)c+")
			Expect(pos).To(Equal(f.Pos(0)))
			Expect(match).To(BeNil())
		})

		It("should match a regexp at the end", func() {
			pos, match := r.ReadRegexpSubmatch(f.Pos(5), "(e)f+")
			Expect(pos).To(Equal(f.Pos(7)))
			Expect(match).To(Equal([][]byte{
				[]byte("ef"),
				[]byte("e"),
			}))
		})

		It("should only match the full match at the end", func() {
			pos, match := r.ReadRegexpSubmatch(f.Pos(5), "efg+")
			Expect(pos).To(Equal(f.Pos(5)))
			Expect(match).To(BeNil())
		})

		Context("at the end of the file", func() {
			It("should not match anything", func() {
				pos, match := r.ReadRegexpSubmatch(f.Pos(7), "x+")
				Expect(pos).To(Equal(f.Pos(7)))
				Expect(match).To(BeNil())
			})
		})

		Context("When input contains UTF8", func() {
			BeforeEach(func() {
				data = []byte(inputWithUTF8)
			})

			It("should match the regexp", func() {
				pos, match := r.ReadRegexpSubmatch(f.Pos(0), "(.*) and")
				Expect(pos).To(Equal(f.Pos(8)))
				Expect(match).To(Equal([][]byte{
					[]byte("🍕 and"),
					[]byte("🍕"),
				}))
			})

			It("should not match a non-matching regexp", func() {
				pos, match := r.ReadRegexpSubmatch(f.Pos(0), ".* not")
				Expect(pos).To(Equal(f.Pos(0)))
				Expect(match).To(BeNil())
			})

			It("should match a regexp at the end", func() {
				pos, match := r.ReadRegexpSubmatch(f.Pos(5), "and (.*)")
				Expect(pos).To(Equal(f.Pos(13)))
				Expect(match).To(Equal([][]byte{
					[]byte("and 🍺"),
					[]byte("🍺"),
				}))
			})

			It("should only match the full match at the end", func() {
				pos, match := r.ReadRegexpSubmatch(f.Pos(5), "and .*s")
				Expect(pos).To(Equal(f.Pos(5)))
				Expect(match).To(BeNil())
			})
		})
	})

	Describe("Readf()", func() {
		var fun func(b []byte) ([]byte, int)

		BeforeEach(func() {
			fun = func(b []byte) ([]byte, int) {
				return b[0:1], 1
			}
		})

		It("should match the input with the function", func() {
			pos, match := r.Readf(f.Pos(0), fun)
			Expect(pos).To(Equal(f.Pos(1)))
			Expect(match).To(Equal([]byte("a")))
		})

		It("should match the input at the end", func() {
			pos, match := r.Readf(f.Pos(6), fun)
			Expect(pos).To(Equal(f.Pos(7)))
			Expect(match).To(Equal([]byte("f")))
		})

		It("should use the returned position instead of the length of the match", func() {
			fun = func(b []byte) ([]byte, int) {
				return b[0:1], 2
			}
			pos, match := r.Readf(f.Pos(1), fun)
			Expect(pos).To(Equal(f.Pos(3)))
			Expect(match).To(Equal([]byte("b")))
		})

		It("should return with no result and unchanged position if no match", func() {
			fun = func(b []byte) ([]byte, int) {
				return nil, 0
			}
			pos, match := r.Readf(f.Pos(2), fun)
			Expect(pos).To(Equal(f.Pos(2)))
			Expect(match).To(BeNil())
		})

		Context("at the end of the file", func() {
			It("should not match anything", func() {
				pos, match := r.Readf(f.Pos(7), fun)
				Expect(pos).To(Equal(f.Pos(7)))
				Expect(match).To(BeNil())
			})
		})

		Context("when the returned position is after the end of file", func() {
			It("should panic", func() {
				fun = func(b []byte) ([]byte, int) {
					return b[0:1], 2
				}
				Expect(func() { r.Readf(f.Pos(6), fun) }).To(Panic())
			})
		})

		Context("when the returned position is before the end of the match", func() {
			It("should panic", func() {
				fun = func(b []byte) ([]byte, int) {
					return b[0:2], 1
				}
				Expect(func() { r.Readf(f.Pos(0), fun) }).To(Panic())
			})
		})

		Context("when the next positon is zero but a match is returned", func() {
			It("should panic", func() {
				fun = func(b []byte) ([]byte, int) {
					return b[0:1], 0
				}
				Expect(func() { r.Readf(f.Pos(0), fun) }).To(Panic())
			})
		})
	})

	Describe("Remaining()", func() {
		It("should return with the remaining bytes", func() {
			Expect(r.Remaining(f.Pos(0))).To(Equal(len(input)))
		})

		It("should return with the remaining bytes from a given position", func() {
			Expect(r.Remaining(f.Pos(3))).To(Equal(len(input) - 3))
		})
	})

	Describe("Pos()", func() {
		It("should return with global pos", func() {
			Expect(r.Pos(1)).To(Equal(parsley.Pos(2)))
		})
	})

	Describe("IsEOF()", func() {
		It("should return false before the end of the input", func() {
			Expect(r.IsEOF(f.Pos(0))).To(BeFalse())
			Expect(r.IsEOF(f.Pos(6))).To(BeFalse())
		})
		It("should return true at the end of the input", func() {
			Expect(r.IsEOF(f.Pos(7))).To(BeTrue())
		})
	})

	Describe("SkipWhitespaces()", func() {
		BeforeEach(func() {
			data = []byte("abc \t\n\fdef  ")
		})

		It("should not match any whitespaces if none", func() {
			pos, err := r.SkipWhitespaces(f.Pos(0), text.WsSpacesNl)
			Expect(pos).To(Equal(f.Pos(0)))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should match all types of whitespaces", func() {
			pos, err := r.SkipWhitespaces(f.Pos(3), text.WsSpacesNl)
			Expect(pos).To(Equal(f.Pos(7)))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when new lines are not valid", func() {
			It("should return an error", func() {
				pos, err := r.SkipWhitespaces(f.Pos(3), text.WsSpaces)
				Expect(pos).To(Equal(f.Pos(7)))
				Expect(err).To(MatchError(parsley.NewError(f.Pos(5), parsley.NewWhitespaceError("new line is not allowed"))))
			})
		})

		Context("when not skipping any whitespaces", func() {
			It("should return an error", func() {
				pos, err := r.SkipWhitespaces(f.Pos(3), text.WsNone)
				Expect(pos).To(Equal(f.Pos(7)))
				Expect(err).To(MatchError(parsley.NewError(f.Pos(3), parsley.NewWhitespaceError("whitespaces are not allowed"))))
			})
		})

		Context("when forcing a new line", func() {
			It("should return an error", func() {
				pos, err := r.SkipWhitespaces(f.Pos(3), text.WsSpacesForceNl)
				Expect(pos).To(Equal(f.Pos(7)))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should not match spaces and tabs only", func() {
				pos, err := r.SkipWhitespaces(f.Pos(10), text.WsSpacesForceNl)
				Expect(pos).To(Equal(f.Pos(12)))
				Expect(err).To(MatchError(parsley.NewError(f.Pos(12), parsley.NewWhitespaceError("was expecting a new line"))))
			})
		})
	})
})
