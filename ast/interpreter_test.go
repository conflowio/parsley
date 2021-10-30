// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ast_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/ast"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/parsley/parsleyfakes"
)

var _ = Describe("Interpreter", func() {

	It("Eval should call the function", func() {
		var (
			passedCtx  interface{}
			passedNode parsley.NonTerminalNode
			fResult    = "some result"
			fErr       = parsley.NewErrorf(parsley.Pos(1), "some error")
			ctx        = "some context"
			node       = &parsleyfakes.FakeNonTerminalNode{}
		)
		f := func(userCtx interface{}, node parsley.NonTerminalNode) (interface{}, parsley.Error) {
			passedCtx = userCtx
			passedNode = node
			return fResult, fErr
		}

		result, err := ast.InterpreterFunc(f).Eval(ctx, node)
		Expect(result).To(Equal(fResult))
		Expect(err).To(Equal(fErr))

		Expect(passedCtx).To(Equal(ctx))
		Expect(passedNode).To(Equal(node))
	})
})
