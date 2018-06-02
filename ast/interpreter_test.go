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

var _ = Describe("Interpreter", func() {

	It("Eval should call the function", func() {
		var (
			passedCtx   interface{}
			passedNodes []parsley.Node
			fResult     = "some result"
			fErr        = parsley.NewErrorf(parsley.Pos(1), "some error")
			ctx         = "some context"
			nodes       = []parsley.Node{nil}
		)
		f := func(ctx interface{}, nodes []parsley.Node) (interface{}, parsley.Error) {
			passedCtx = ctx
			passedNodes = nodes
			return fResult, fErr
		}

		result, err := ast.InterpreterFunc(f).Eval(ctx, nodes)
		Expect(result).To(Equal(fResult))
		Expect(err).To(Equal(fErr))

		Expect(passedCtx).To(Equal(ctx))
		Expect(passedNodes).To(Equal(nodes))
	})
})
