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
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Float", func() {

	var p = terminal.Float()

	DescribeTable("should match",
		func(input string, startPos int, value interface{}, nodePos parsley.Pos, endPos int) {
			r := text.NewReader(text.NewFile("textfile", []byte(input)))
			curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			node := res[0].(*ast.TerminalNode)
			Expect(node.Token()).To(Equal("FLOAT"))
			Expect(node.Value(nil)).To(Equal(value))
			Expect(node.Pos()).To(Equal(nodePos))
			Expect(node.ReaderPos()).To(Equal(endPos))
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
		func(input string, startPos int, errPos parsley.Pos) {
			r := text.NewReader(text.NewFile("textfile", []byte(input)))
			curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).To(MatchError("was expecting float value"))
			Expect(err.Pos()).To(Equal(errPos))
			Expect(res).To(BeNil())
		},
		Entry("empty", "", 0, parsley.Pos(1)),
		Entry("pos test", "--- x", 4, parsley.Pos(5)),

		Entry("a", "a", 0, parsley.Pos(1)),
		Entry("-", "-", 0, parsley.Pos(1)),
		Entry("+", "+", 0, parsley.Pos(1)),
		Entry(".", ".", 0, parsley.Pos(1)),
		Entry("1", "1", 0, parsley.Pos(1)),
		Entry("-1", "-1", 0, parsley.Pos(1)),
		Entry("0", "0", 0, parsley.Pos(1)),
	)

	Describe("when there is an invalid float value", func() {
		It("should trow an error", func() {
			input := "1.2e3456"
			r := text.NewReader(text.NewFile("textfile", []byte(input)))
			curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, 0)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).To(MatchError("invalid float value encountered"))
			Expect(err.Pos()).To(Equal(parsley.Pos(1)))
			Expect(res).To(BeNil())
		})
	})
})
