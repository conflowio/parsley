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

var _ = Describe("NL", func() {

	var p = terminal.NL()

	DescribeTable("should match",
		func(input string, startPos int, value interface{}, nodePos int, endPos int) {
			f := text.NewFile("textfile", []byte(input))
			fs := parsley.NewFileSet(f)
			r := text.NewReader(f)
			ctx := parsley.NewContext(fs, r)
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, parsley.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Token()).To(Equal("NL"))
			Expect(res.Value(nil)).To(Equal(value))
			Expect(res.Pos()).To(Equal(parsley.Pos(nodePos)))
			Expect(res.ReaderPos()).To(Equal(parsley.Pos(endPos)))
		},
		Entry(`single nl`, "\n", 1, "\n", 1, 2),
		Entry(`single nl, ws before`, " \t\n", 1, "\n", 3, 4),
		Entry(`single nl, ws before/after`, " \t\n \t", 1, "\n", 3, 4),
		Entry(`two nls`, " \t\n \t\n", 1, "\n\n", 3, 7),
		Entry(`three nls`, "\n\n\n", 1, "\n\n", 1, 4),
	)

	DescribeTable("should not match",
		func(input string, startPos int, errPos int) {
			f := text.NewFile("textfile", []byte(input))
			fs := parsley.NewFileSet(f)
			r := text.NewReader(f)
			ctx := parsley.NewContext(fs, r)
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, parsley.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).To(HaveOccurred())
			Expect(err.Cause()).To(MatchError("was expecting new line"))
			Expect(err.Pos()).To(Equal(parsley.Pos(errPos)))
			Expect(res).To(BeNil())
		},
		Entry("empty", ``, 1, 1),
		Entry("not nl", `x`, 1, 1),
		Entry("spaces only", `  `, 1, 3),
	)
})
