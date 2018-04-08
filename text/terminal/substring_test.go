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

var _ = Describe("Substring", func() {

	var p = terminal.Substring("FOO", "foo", 42, text.WsNone)

	It("should have a name", func() {
		Expect(p.Name()).To(Equal(`"foo"`))
	})

	Context("when called with an empty string", func() {
		It("should panic", func() {
			Expect(func() { terminal.Substring("FOO", "", 42, text.WsNone) }).To(Panic())
		})
	})

	DescribeTable("should match",
		func(input string, startPos int, value interface{}, nodePos parsley.Pos, endPos int) {
			r := text.NewReader(text.NewFile("textfile", []byte(input)))
			curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			node := res.(*ast.TerminalNode)
			Expect(node.Token()).To(Equal("FOO"))
			Expect(node.Value(nil)).To(Equal(value))
			Expect(node.Pos()).To(Equal(nodePos))
			Expect(node.ReaderPos()).To(Equal(endPos))
		},
		Entry(`foo beginning`, `foo`, 0, 42, parsley.Pos(1), 3),
		Entry(`foo middle`, `--- foo ---`, 4, 42, parsley.Pos(5), 7),
		Entry(`foo end`, `--- foo`, 4, 42, parsley.Pos(5), 7),
		Entry(`prefix`, `foobar`, 0, 42, parsley.Pos(1), 3),
	)

	DescribeTable("should not match",
		func(input string, startPos int) {
			r := text.NewReader(text.NewFile("textfile", []byte(input)))
			curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(BeNil())
		},
		Entry("empty", ``, 0),
		Entry("partial", `fo`, 0),
	)

	Context("with wsSpaces", func() {

		var p = terminal.Substring("FOO", "foo", nil, text.WsSpaces)

		It("should skip spaces and tabs", func() {
			r := text.NewReader(text.NewFile("textfile", []byte("foo \t\n\fxxx")))
			_, res, _ := p.Parse(nil, data.EmptyIntMap, r, 0)
			Expect(res.ReaderPos()).To(Equal(5))
		})
	})

	Context("with wsSpacesNl", func() {

		var p = terminal.Substring("FOO", "foo", nil, text.WsSpacesNl)

		It("should skip spaces, tabs and new lines", func() {
			r := text.NewReader(text.NewFile("textfile", []byte("foo \t\n\fxxx")))
			_, res, _ := p.Parse(nil, data.EmptyIntMap, r, 0)
			Expect(res.ReaderPos()).To(Equal(7))
		})
	})
})
