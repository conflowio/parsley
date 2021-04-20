// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("EmptyNode", func() {
	var (
		node ast.EmptyNode
		pos  parsley.Pos = parsley.Pos(1)
	)

	JustBeforeEach(func() {
		node = ast.EmptyNode(pos)
	})

	Describe("Methods", func() {
		It("Token() should return with the token value", func() {
			Expect(node.Token()).To(Equal("EMPTY"))
		})

		It("Pos() should return with the token value", func() {
			Expect(node.Pos()).To(Equal(pos))
		})

		It("ReaderPos() should return with the reader position", func() {
			Expect(node.ReaderPos()).To(Equal(pos))
		})

		It("String() should return with NIL", func() {
			Expect(node.String()).To(Equal("EMPTY{1}"))
		})
	})
})
