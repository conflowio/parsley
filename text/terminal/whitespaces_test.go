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

var _ = Describe("Whitespaces", func() {

	Context("when not allowing new lines", func() {

		var p = terminal.Whitespaces(false)

		DescribeTable("should match",
			func(input string, startPos int, nodePos parsley.Pos, endPos int) {
				r := text.NewReader(text.NewFile("textfile", []byte(input)))
				curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
				Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
				Expect(err).ToNot(HaveOccurred())
				node := res[0].(*ast.TerminalNode)
				Expect(node.Token()).To(Equal("WS"))
				Expect(node.Value(nil)).To(BeNil())
				Expect(node.Pos()).To(Equal(nodePos))
				Expect(node.ReaderPos()).To(Equal(endPos))
			},
			Entry("ws beginning", " \t---", 0, parsley.Pos(1), 2),
			Entry("ws middle", "--- \t---", 3, parsley.Pos(4), 5),
			Entry("ws end", "--- \t", 3, parsley.Pos(4), 5),
			Entry("should not match new line", " \t\n\f", 0, parsley.Pos(1), 2),
		)

		DescribeTable("should not match",
			func(input string, startPos int, errPos parsley.Pos) {
				r := text.NewReader(text.NewFile("textfile", []byte(input)))
				curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
				Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
				Expect(err).To(MatchError("was expecting one or more whitespaces"))
				Expect(err.Pos()).To(Equal(errPos))
				Expect(res).To(BeNil())
			},
			Entry("empty", "", 0, parsley.Pos(1)),
			Entry("pos test", `--- x`, 4, parsley.Pos(5)),
			Entry("not whitespace", `a`, 0, parsley.Pos(1)),
			Entry("new line", `\n\t`, 0, parsley.Pos(1)),
		)
	})

	Context("when allowing new lines", func() {

		var p = terminal.Whitespaces(true)

		DescribeTable("should match (with new lines)",
			func(input string, startPos int, nodePos parsley.Pos, endPos int) {
				r := text.NewReader(text.NewFile("textfile", []byte(input)))
				curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
				Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
				Expect(err).ToNot(HaveOccurred())
				node := res[0].(*ast.TerminalNode)
				Expect(node.Token()).To(Equal("WS"))
				Expect(node.Value(nil)).To(BeNil())
				Expect(node.Pos()).To(Equal(nodePos))
				Expect(node.ReaderPos()).To(Equal(endPos))
			},
			Entry("ws beginning", " \t\n\f---", 0, parsley.Pos(1), 4),
			Entry("ws middle", "--- \t\n\f---", 3, parsley.Pos(4), 7),
			Entry("ws end", "--- \t\n\f", 3, parsley.Pos(4), 7),
		)

		DescribeTable("should not match (with new lines)",
			func(input string, startPos int, errPos parsley.Pos) {
				r := text.NewReader(text.NewFile("textfile", []byte(input)))
				curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
				Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
				Expect(err).To(MatchError("was expecting one or more whitespaces"))
				Expect(err.Pos()).To(Equal(errPos))
				Expect(res).To(BeNil())
			},
			Entry("empty", "", 0, parsley.Pos(1)),
			Entry("pos test", `--- x`, 4, parsley.Pos(5)),
			Entry("not whitespace", `a`, 0, parsley.Pos(1)),
		)
	})

})
