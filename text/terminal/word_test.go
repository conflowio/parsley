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

var _ = Describe("Word", func() {

	var p = terminal.Word("integer", "foo", 42)

	Context("when called with an empty word", func() {
		It("should panic", func() {
			Expect(func() { terminal.Word("integer", "", 42) }).To(Panic())
		})
	})

	DescribeTable("should match",
		func(input string, startPos int, value interface{}, nodePos parsley.Pos, endPos int) {
			f := text.NewFile("textfile", []byte(input))
			fs := parsley.NewFileSet(f)
			r := text.NewReader(f)
			ctx := parsley.NewContext(fs, r)
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Token()).To(Equal("FOO"))
			Expect(res.Schema()).To(Equal("integer"))
			Expect(res.(parsley.LiteralNode).Value()).To(Equal(value))
			Expect(res.Pos()).To(Equal(nodePos))
			Expect(res.ReaderPos()).To(Equal(f.Pos(endPos)))
		},
		Entry(`foo beginning`, `foo`, 0, 42, parsley.Pos(1), 3),
		Entry(`foo middle`, `--- foo ---`, 4, 42, parsley.Pos(5), 7),
		Entry(`foo end`, `--- foo`, 4, 42, parsley.Pos(5), 7),
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
			Expect(err.Cause()).To(MatchError("was expecting \"foo\""))
			Expect(err.Pos()).To(Equal(f.Pos(startPos)))
			Expect(res).To(BeNil())
		},
		Entry("empty", ``, 0),
		Entry("prefix", `foobar`, 0),
		Entry("partial", `fo`, 0),
	)
})
