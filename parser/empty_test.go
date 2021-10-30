// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/parsley/parsleyfakes"
)

var _ = Describe("Empty", func() {

	var (
		p   = parser.Empty()
		ctx *parsley.Context
	)

	BeforeEach(func() {
		ctx = parsley.NewContext(parsley.NewFileSet(), &parsleyfakes.FakeReader{})
	})

	It("should return with an empty node", func() {
		res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, 1)
		Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
		Expect(res).To(Equal(ast.EmptyNode(1)))
		Expect(err).ToNot(HaveOccurred())
	})

})
