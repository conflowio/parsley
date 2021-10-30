// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/parsley/parsleyfakes"
)

var _ = Describe("AppendNode", func() {

	var (
		n1, n2, n3, n4 *parsleyfakes.FakeNode
	)

	BeforeEach(func() {
		n1 = &parsleyfakes.FakeNode{}
		n1.TokenReturns("n1")
		n2 = &parsleyfakes.FakeNode{}
		n2.TokenReturns("n2")
		n3 = &parsleyfakes.FakeNode{}
		n3.TokenReturns("n3")
		n4 = &parsleyfakes.FakeNode{}
		n4.TokenReturns("n4")
	})

	DescribeTable("result",
		func(n1 parsley.Node, n2 parsley.Node, res parsley.Node) {
			if res != nil {
				Expect(ast.AppendNode(n1, n2)).To(Equal(res))
			} else {
				Expect(ast.AppendNode(n1, n2)).To(BeNil())
			}
		},
		Entry("nil + nil", nil, nil, nil),
		Entry("node + nil", n1, nil, n1),
		Entry("nil + node", nil, n1, n1),
		Entry("node + node", n1, n2, ast.NodeList([]parsley.Node{n1, n2})),
		Entry("node + list", n1, ast.NodeList([]parsley.Node{n2, n3}), ast.NodeList([]parsley.Node{n1, n2, n3})),
		Entry("list + node", ast.NodeList([]parsley.Node{n1, n2}), n3, ast.NodeList([]parsley.Node{n1, n2, n3})),
		Entry(
			"list + list",
			ast.NodeList([]parsley.Node{n1, n2}),
			ast.NodeList([]parsley.Node{n3, n4}),
			ast.NodeList([]parsley.Node{n1, n2, n3, n4}),
		),
	)
})

var _ = Describe("SetReaderPos", func() {
	It("calls the SetReaderPos function if node implements the SetReaderPos interface", func() {
		node := ast.NewTerminalNode("string", "TEST", "x", parsley.Pos(1), parsley.Pos(2))
		f := func(pos parsley.Pos) parsley.Pos {
			return parsley.Pos(pos + 1)
		}

		resultNode := ast.SetReaderPos(node, f)

		Expect(resultNode).To(BeEquivalentTo(node))
		Expect(node.ReaderPos()).To(Equal(parsley.Pos(3)))
	})

	It("returns a new NilPos with the new position", func() {
		node := ast.EmptyNode(parsley.Pos(1))
		f := func(pos parsley.Pos) parsley.Pos {
			return parsley.Pos(pos + 1)
		}

		resultNode := ast.SetReaderPos(node, f)

		Expect(resultNode.ReaderPos()).To(Equal(parsley.Pos(2)))
	})

	It("sets the reader position for all nodes in a node list", func() {
		n1 := ast.NewTerminalNode("string", "TEST", "x", parsley.Pos(1), parsley.Pos(2))
		n2 := ast.EmptyNode(parsley.Pos(3))

		nl := ast.NodeList([]parsley.Node{n1, n2})

		f := func(pos parsley.Pos) parsley.Pos {
			return parsley.Pos(pos + 1)
		}

		resultNode := ast.SetReaderPos(nl, f)

		Expect(resultNode).To(BeEquivalentTo(nl))
		Expect(nl[0].ReaderPos()).To(Equal(parsley.Pos(3)))
		Expect(nl[1].ReaderPos()).To(Equal(parsley.Pos(4)))
	})

	It("panics for an unknown node type", func() {
		node := &parsleyfakes.FakeNode{}

		f := func(pos parsley.Pos) parsley.Pos {
			return parsley.Pos(pos + 1)
		}

		Expect(func() { ast.SetReaderPos(node, f) }).To(Panic())
	})
})
