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
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
	"github.com/opsidian/parsley/text/terminal"
)

var _ = Describe("Nil", func() {

	var p = terminal.Nil("nil", "nil")

	Context("when called with an empty nil value", func() {
		It("should panic", func() {
			Expect(func() { terminal.Nil("nil", "") }).To(Panic())
		})
	})

	Context("when called with an empty nil type", func() {
		It("should panic", func() {
			Expect(func() { terminal.Nil("nil", "") }).To(Panic())
		})
	})

	DescribeTable("should match",
		func(input string, startPos int, nodePos parsley.Pos, endPos int) {
			f := text.NewFile("textfile", []byte(input))
			fs := parsley.NewFileSet(f)
			r := text.NewReader(f)
			ctx := parsley.NewContext(fs, r)
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Token()).To(Equal("NIL"))
			Expect(res.(parsley.LiteralNode).Value()).To(BeNil())
			Expect(res.Pos()).To(Equal(nodePos))
			Expect(res.ReaderPos()).To(Equal(f.Pos(endPos)))
		},
		Entry("nil beginning", "nil ---", 0, parsley.Pos(1), 3),
		Entry("nil middle", "--- nil ---", 4, parsley.Pos(5), 7),
		Entry("nil end", "--- nil", 4, parsley.Pos(5), 7),
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
			Expect(err.Cause()).To(MatchError("was expecting nil"))
			Expect(err.Pos()).To(Equal(f.Pos(startPos)))
			Expect(res).To(BeNil())
		},
		Entry("empty", "", 0),
		Entry("nilx", "truex", 0),
		Entry("ni", "ni", 0),
		Entry("Nil", "Nil", 0),
	)

})
