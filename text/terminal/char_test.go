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

var _ = Describe("Char", func() {

	var p = terminal.Char()

	It("should have a name", func() {
		Expect(p.Name()).ToNot(BeEmpty())
	})

	DescribeTable("should match",
		func(input string, startPos int, value interface{}, nodePos parsley.Pos, endPos int) {
			f := text.NewFile("textfile", []byte(input))
			r := text.NewReader(f)
			res, err, curtailingParsers := p.Parse(nil, data.EmptyIntMap, r, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			node := res.(*ast.TerminalNode)
			Expect(node.Token()).To(Equal("CHAR"))
			Expect(node.Value(nil)).To(Equal(value))
			Expect(node.Pos()).To(Equal(nodePos))
			Expect(node.ReaderPos()).To(Equal(f.Pos(endPos)))
		},
		Entry(`'a' beginning`, `'a'`, 0, 'a', parsley.Pos(1), 3),
		Entry(`'a' middle`, `--- 'a' ---`, 4, 'a', parsley.Pos(5), 7),
		Entry(`'a' end`, `--- 'a'`, 4, 'a', parsley.Pos(5), 7),

		Entry(`quote`, `'\''`, 0, '\'', parsley.Pos(1), 4),
		Entry(`double quote`, `'"'`, 0, '"', parsley.Pos(1), 3),

		Entry(`' '`, `' '`, 0, ' ', parsley.Pos(1), 3),
		Entry(`'🍕'`, `'🍕'`, 0, '🍕', parsley.Pos(1), 6),
		Entry(`'\a'`, `'\a'`, 0, '\a', parsley.Pos(1), 4),
		Entry(`'\b'`, `'\b'`, 0, '\b', parsley.Pos(1), 4),
		Entry(`'\f'`, `'\f'`, 0, '\f', parsley.Pos(1), 4),
		Entry(`'\n'`, `'\n'`, 0, '\n', parsley.Pos(1), 4),
		Entry(`'\v'`, `'\v'`, 0, '\v', parsley.Pos(1), 4),
		Entry(`'\r'`, `'\r'`, 0, '\r', parsley.Pos(1), 4),
		Entry(`'\t'`, `'\t'`, 0, '\t', parsley.Pos(1), 4),
		Entry(`'\xa0'`, `'\xa0'`, 0, '\xa0', parsley.Pos(1), 6),
		Entry(`'\xA0'`, `'\xA0'`, 0, '\xa0', parsley.Pos(1), 6),
		Entry(`'\u12ab'`, `'\u12ab'`, 0, '\u12ab', parsley.Pos(1), 8),
		Entry(`'\u12AB'`, `'\u12AB'`, 0, '\u12ab', parsley.Pos(1), 8),
		Entry(`'\U0001f355'`, `'\U0001f355'`, 0, '🍕', parsley.Pos(1), 12),
		Entry(`'\U0001F355'`, `'\U0001F355'`, 0, '🍕', parsley.Pos(1), 12),
	)

	DescribeTable("should not match",
		func(input string, startPos int) {
			f := text.NewFile("textfile", []byte(input))
			r := text.NewReader(f)
			res, err, curtailingParsers := p.Parse(nil, data.EmptyIntMap, r, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			Expect(res).To(BeNil())
		},
		Entry("empty", ``, 0),
		Entry(`x`, `x`, 0),
		Entry(`"a"`, `"a"`, 0),
		Entry(`\x`, `\x`, 0),
		Entry(`\u`, `\u`, 0),
	)

	DescribeTable("should error",
		func(input string, startPos int, errPos parsley.Pos) {
			f := text.NewFile("textfile", []byte(input))
			r := text.NewReader(f)
			res, err, curtailingParsers := p.Parse(nil, data.EmptyIntMap, r, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())
		},
		Entry("only start quote", `'`, 0, parsley.Pos(1)),
		Entry("empty quotes", `''`, 0, parsley.Pos(2)),
		Entry("no end quote", `'a`, 0, parsley.Pos(3)),
		Entry("multiple characters", `'aa'`, 0, parsley.Pos(3)),
	)
})
