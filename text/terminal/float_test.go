// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	"github.com/conflowio/parsley/text/terminal"
)

var _ = Describe("Float", func() {

	var p = terminal.Float("float")

	DescribeTable("should match",
		func(input string, startPos int, value interface{}, nodePos parsley.Pos, endPos int) {
			f := text.NewFile("textfile", []byte(input))
			fs := parsley.NewFileSet(f)
			r := text.NewReader(f)
			ctx := parsley.NewContext(fs, r)
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Token()).To(Equal("FLOAT"))
			Expect(res.Schema()).To(Equal("float"))
			Expect(res.(parsley.LiteralNode).Value()).To(Equal(value))
			Expect(res.Pos()).To(Equal(nodePos))
			Expect(res.ReaderPos()).To(Equal(f.Pos(endPos)))
		},
		Entry("1.2 beginning", "1.2 ---", 0, 1.2, parsley.Pos(1), 3),
		Entry("1.2 middle", "--- 1.2 ---", 4, 1.2, parsley.Pos(5), 7),
		Entry("1.2 end", "--- 1.2", 4, 1.2, parsley.Pos(5), 7),
		Entry("0.1", "0.1", 0, 0.1, parsley.Pos(1), 3),
		Entry("+0.1", "+0.1", 0, 0.1, parsley.Pos(1), 4),
		Entry("-0.1", "-0.1", 0, -0.1, parsley.Pos(1), 4),
		Entry("1.1", "1.1", 0, 1.1, parsley.Pos(1), 3),
		Entry("+1.1", "+1.1", 0, 1.1, parsley.Pos(1), 4),
		Entry("-1.1", "-1.1", 0, -1.1, parsley.Pos(1), 4),
		Entry("1234567890.1234567890", "1234567890.1234567890", 0, 1234567890.1234567890, parsley.Pos(1), 21),
		Entry("+1234567890.1234567890", "+1234567890.1234567890", 0, 1234567890.1234567890, parsley.Pos(1), 22),
		Entry("-1234567890.1234567890", "-1234567890.1234567890", 0, -1234567890.1234567890, parsley.Pos(1), 22),
		Entry("123.4abc", "123.4abc", 0, 123.4, parsley.Pos(1), 5),
		Entry("00.00", "00.00", 0, 0.0, parsley.Pos(1), 5),
		Entry(".0", ".0", 0, 0.0, parsley.Pos(1), 2),
		Entry(".1", ".1", 0, 0.1, parsley.Pos(1), 2),
		Entry("+.1", "+.1", 0, 0.1, parsley.Pos(1), 3),
		Entry("-.1", "-.1", 0, -0.1, parsley.Pos(1), 3),
		Entry("1.2e5", "1.2e5", 0, 1.2e5, parsley.Pos(1), 5),
		Entry("+1.2e5", "+1.2e5", 0, 1.2e5, parsley.Pos(1), 6),
		Entry("-1.2e5", "-1.2e5", 0, -1.2e5, parsley.Pos(1), 6),
		Entry("1.2e", "1.2e", 0, 1.2, parsley.Pos(1), 3), // only 1.2 should be consumed
	)

	DescribeTable("should not match",
		func(input string, startPos int) {
			f := text.NewFile("textfile", []byte(input))
			fs := parsley.NewFileSet(f)
			r := text.NewReader(f)
			ctx := parsley.NewContext(fs, r)
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).To(HaveOccurred())
			Expect(err.Cause()).To(MatchError("was expecting float value"))
			Expect(err.Pos()).To(Equal(f.Pos(startPos)))
			Expect(res).To(BeNil())
		},
		Entry("empty", "", 0),
		Entry("a", "a", 0),
		Entry("-", "-", 0),
		Entry("+", "+", 0),
		Entry(".", ".", 0),
		Entry("1", "1", 0),
		Entry("-1", "-1", 0),
		Entry("0", "0", 0),
	)

	Describe("when there is an invalid float value", func() {
		It("should trow an error", func() {
			input := "1.2e3456"
			f := text.NewFile("textfile", []byte(input))
			fs := parsley.NewFileSet(f)
			r := text.NewReader(f)
			ctx := parsley.NewContext(fs, r)
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, f.Pos(0))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).To(MatchError("invalid float value"))
			Expect(err.Pos()).To(Equal(parsley.Pos(1)))
			Expect(res).To(BeNil())
		})
	})
})
