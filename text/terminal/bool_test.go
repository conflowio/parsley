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

var _ = Describe("Bool", func() {

	var p = terminal.Bool("true", "false")

	It("should have a name", func() {
		Expect(p.Name()).To(Equal("true or false"))
	})

	Context("when called with an empty true/false value", func() {
		It("should panic", func() {
			Expect(func() { terminal.Bool("", "false") }).To(Panic())
			Expect(func() { terminal.Bool("true", "") }).To(Panic())
		})
	})

	DescribeTable("should match",
		func(input string, startPos int, value interface{}, nodePos parsley.Pos, endPos int) {
			f := text.NewFile("textfile", []byte(input))
			r := text.NewReader(f)
			curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			node := res.(*ast.TerminalNode)
			Expect(node.Token()).To(Equal("BOOL"))
			Expect(node.Value(nil)).To(Equal(value))
			Expect(node.Pos()).To(Equal(nodePos))
			Expect(node.ReaderPos()).To(Equal(f.Pos(endPos)))
		},
		Entry("true beginning", "true ---", 0, true, parsley.Pos(1), 4),
		Entry("false beginning", "false ---", 0, false, parsley.Pos(1), 5),
		Entry("true middle", "--- true ---", 4, true, parsley.Pos(5), 8),
		Entry("false middle", "--- false ---", 4, false, parsley.Pos(5), 9),
		Entry("true end", "--- true", 4, true, parsley.Pos(5), 8),
		Entry("false end", "--- false", 4, false, parsley.Pos(5), 9),
	)

	DescribeTable("should not match",
		func(input string, startPos int) {
			f := text.NewFile("textfile", []byte(input))
			r := text.NewReader(f)
			curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(BeNil())
		},
		Entry("empty", "", 0),
		Entry("truex", "truex", 0),
		Entry("falsex", "falsex", 0),
		Entry("tru", "tru", 0),
		Entry("fals", "fals", 0),
		Entry("True", "True", 0),
		Entry("False", "False", 0),
	)

})
