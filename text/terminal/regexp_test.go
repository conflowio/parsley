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

var _ = Describe("Regexp", func() {

	var p1 = terminal.Regexp("FOO", "foo", "fo+", 0)
	var p2 = terminal.Regexp("FOO", "foo", "f(o+)", 1)

	It("should have a name", func() {
		Expect(p1.Name()).To(Equal("foo"))
	})

	Context("when regexp matches an empty string", func() {
		It("should panic", func() {
			r := text.NewReader(text.NewFile("textfile", []byte("foo")))
			ctx := parsley.NewContext(r)
			p := terminal.Regexp("FOO", "foo", "f*", 0)
			Expect(func() { p.Parse(ctx, data.EmptyIntMap, 0) }).To(Panic())
		})
	})

	Context("when capturing group is invalid", func() {
		It("should panic", func() {
			r := text.NewReader(text.NewFile("textfile", []byte("foo")))
			ctx := parsley.NewContext(r)
			p := terminal.Regexp("FOO", "foo", "f(o+)", 2)
			Expect(func() { p.Parse(ctx, data.EmptyIntMap, 0) }).To(Panic())
		})
	})

	DescribeTable("full match - should match",
		func(input string, startPos int, value interface{}, nodePos parsley.Pos, endPos int) {
			f := text.NewFile("textfile", []byte(input))
			r := text.NewReader(f)
			ctx := parsley.NewContext(r)
			res, curtailingParsers := p1.Parse(ctx, data.EmptyIntMap, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(ctx.Error()).ToNot(HaveOccurred())
			node := res.(*ast.TerminalNode)
			Expect(node.Token()).To(Equal("FOO"))
			Expect(node.Value(nil)).To(Equal(value))
			Expect(node.Pos()).To(Equal(nodePos))
			Expect(node.ReaderPos()).To(Equal(f.Pos(endPos)))
		},
		Entry(`foo beginning`, `foo ---`, 0, "foo", parsley.Pos(1), 3),
		Entry(`foo middle`, `--- foo ---`, 4, "foo", parsley.Pos(5), 7),
		Entry(`foo end`, `--- foo`, 4, "foo", parsley.Pos(5), 7),
	)

	DescribeTable("full match - should not match",
		func(input string, startPos int) {
			f := text.NewFile("textfile", []byte(input))
			r := text.NewReader(f)
			ctx := parsley.NewContext(r)
			res, curtailingParsers := p1.Parse(ctx, data.EmptyIntMap, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(ctx.Error()).ToNot(HaveOccurred())
			Expect(res).To(BeNil())
		},
		Entry("empty", ``, 0),
		Entry("other", `bar`, 0),
	)

	DescribeTable("submatch - should match",
		func(input string, startPos int, value interface{}, nodePos parsley.Pos, endPos int) {
			f := text.NewFile("textfile", []byte(input))
			r := text.NewReader(f)
			ctx := parsley.NewContext(r)
			res, curtailingParsers := p2.Parse(ctx, data.EmptyIntMap, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(ctx.Error()).ToNot(HaveOccurred())
			node := res.(*ast.TerminalNode)
			Expect(node.Token()).To(Equal("FOO"))
			Expect(node.Value(nil)).To(Equal(value))
			Expect(node.Pos()).To(Equal(nodePos))
			Expect(node.ReaderPos()).To(Equal(f.Pos(endPos)))
		},
		Entry(`foo beginning`, `foo ---`, 0, "oo", parsley.Pos(1), 3),
		Entry(`foo middle`, `--- foo ---`, 4, "oo", parsley.Pos(5), 7),
		Entry(`foo end`, `--- foo`, 4, "oo", parsley.Pos(5), 7),
	)

	DescribeTable("submatch - should not match",
		func(input string, startPos int) {
			f := text.NewFile("textfile", []byte(input))
			r := text.NewReader(f)
			ctx := parsley.NewContext(r)
			res, curtailingParsers := p2.Parse(ctx, data.EmptyIntMap, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(ctx.Error()).ToNot(HaveOccurred())
			Expect(res).To(BeNil())
		},
		Entry("empty", ``, 0),
		Entry("other", `bar`, 0),
	)
})
