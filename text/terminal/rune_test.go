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

var _ = Describe("Rune", func() {

	var p = terminal.Rune('+')

	It("should have a name", func() {
		Expect(p.Name()).To(Equal(`"+"`))
	})

	DescribeTable("should match",
		func(input string, startPos int, value interface{}, nodePos parsley.Pos, endPos int) {
			r := text.NewReader(text.NewFile("textfile", []byte(input)))
			curtailingParsers, res, err := p.Parse(nil, data.EmptyIntMap, r, startPos)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			node := res[0].(*ast.TerminalNode)
			Expect(node.Token()).To(Equal("+"))
			Expect(node.Value(nil)).To(Equal(value))
			Expect(node.Pos()).To(Equal(nodePos))
			Expect(node.ReaderPos()).To(Equal(endPos))
		},
		Entry(`+ beginning`, `+ ---`, 0, '+', parsley.Pos(1), 1),
		Entry(`+ middle`, `--- + ---`, 4, '+', parsley.Pos(5), 5),
		Entry(`+ end`, `--- +`, 4, '+', parsley.Pos(5), 5),
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
		Entry("x", `x`, 0),
	)
})
