// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("End", func() {

	var (
		p   = parser.End()
		r   *parsleyfakes.FakeReader
		ctx *parsley.Context
	)

	BeforeEach(func() {
		r = &parsleyfakes.FakeReader{}
		ctx = parsley.NewContext(parsley.NewFileSet(), r)
	})

	Context("when at the end of the input", func() {
		It("should return with an EOF node", func() {
			r.IsEOFReturns(true)
			r.PosReturns(parsley.Pos(2))
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, 2)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(res).To(Equal(parser.EndNode(parsley.Pos(2))))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when not at the end of the input", func() {
		It("should return with a nil result", func() {
			r.IsEOFReturns(false)
			r.PosReturns(parsley.Pos(2))
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, 1)
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("was expecting the end of input"))
			Expect(err.Pos()).To(Equal(parsley.Pos(1)))
		})
	})

})
