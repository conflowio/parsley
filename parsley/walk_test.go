// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/parsley/parsleyfakes"
)

var _ = Describe("Walk", func() {
	It("calls the Walk function if node implements the Walkable interface", func() {
		n1 := &parsleyfakes.FakeNode{}
		n1.TokenReturns("N1")
		n2 := &parsleyfakes.FakeNode{}
		n2.TokenReturns("N2")

		n := &parsleyfakes.FakeWalkableNode{}
		n.WalkReturns(true)

		f := func(n parsley.Node) bool {
			return false
		}

		res := parsley.Walk(n, f)
		Expect(res).To(BeTrue())

		Expect(n.WalkCallCount()).To(Equal(1))
	})

	Context("when the node is a nonterminal node", func() {
		var n1, n2 *parsleyfakes.FakeNode
		var n *parsleyfakes.FakeNonTerminalNode

		BeforeEach(func() {
			n1 = &parsleyfakes.FakeNode{}
			n1.TokenReturns("N1")
			n2 = &parsleyfakes.FakeNode{}
			n2.TokenReturns("N2")

			n = &parsleyfakes.FakeNonTerminalNode{}
		})

		It("calls the Walk function on the children, then on the node", func() {

			n.ChildrenReturns([]parsley.Node{n1, n2})

			called := []parsley.Node{}
			f := func(n parsley.Node) bool {
				called = append(called, n)
				return false
			}

			res := parsley.Walk(n, f)
			Expect(res).To(BeFalse())

			Expect(called).To(Equal([]parsley.Node{n1, n2, n}))
		})

		It("stops if the function returns true", func() {

			n.ChildrenReturns([]parsley.Node{n1, n2})

			called := []parsley.Node{}
			f := func(n parsley.Node) bool {
				called = append(called, n)
				return true
			}

			res := parsley.Walk(n, f)
			Expect(res).To(BeTrue())

			Expect(called).To(Equal([]parsley.Node{n1}))
		})
	})

	It("calls the function for a non-walkable node", func() {
		node := &parsleyfakes.FakeNode{}
		node.TokenReturns("TEST")
		called := []parsley.Node{}
		f := func(n parsley.Node) bool {
			called = append(called, n)
			return false
		}

		res := parsley.Walk(node, f)
		Expect(res).To(BeFalse())

		Expect(called).To(Equal([]parsley.Node{node}))
	})
})
