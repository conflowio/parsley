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

var _ = Describe("Bool", func() {

	var p = terminal.Bool("bool", "true", "false")

	Context("when called with an empty true/false value", func() {
		It("should panic", func() {
			Expect(func() { terminal.Bool("bool", "", "false") }).To(Panic())
			Expect(func() { terminal.Bool("bool", "true", "") }).To(Panic())
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
			Expect(res.Token()).To(Equal("BOOL"))
			Expect(res.Schema()).To(Equal("bool"))
			Expect(res.Value(nil)).To(Equal(value))
			Expect(res.Pos()).To(Equal(nodePos))
			Expect(res.ReaderPos()).To(Equal(f.Pos(endPos)))
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
			fs := parsley.NewFileSet(f)
			r := text.NewReader(f)
			ctx := parsley.NewContext(fs, r)
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).To(HaveOccurred())
			Expect(err.Cause()).To(MatchError("was expecting boolean"))
			Expect(err.Pos()).To(Equal(f.Pos(startPos)))
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
