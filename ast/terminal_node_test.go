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

var _ = Describe("TerminalNode", func() {
	var (
		node      *ast.TerminalNode
		token     = "TEST"
		value     = "some value"
		schema    = "string"
		pos       = parsley.Pos(1)
		readerPos = parsley.Pos(2)
	)

	JustBeforeEach(func() {
		node = ast.NewTerminalNode(schema, token, value, pos, readerPos)
	})

	Describe("Methods", func() {
		It("Token() should return with the token value", func() {
			Expect(node.Token()).To(Equal(token))
		})

		It("Schema() should return with the value's type", func() {
			Expect(node.Schema()).To(Equal(schema))
		})

		It("EvaluateNode() should return with the value", func() {
			nodeValue := node.Value()
			Expect(nodeValue).To(Equal(value))
		})

		It("Pos() should return with the token value", func() {
			Expect(node.Pos()).To(Equal(pos))
		})

		It("ReaderPos() should return with the reader position", func() {
			Expect(node.ReaderPos()).To(Equal(readerPos))
		})

		It("SetReaderPos() should modify the reader position", func() {
			node.SetReaderPos(func(pos parsley.Pos) parsley.Pos {
				return parsley.Pos(pos + 1)
			})
			Expect(node.ReaderPos()).To(Equal(parsley.Pos(3)))
		})

		It("String() should return with a readable representation", func() {
			Expect(node.String()).To(Equal("TEST{some value, 1..2}"))
		})
	})
})
