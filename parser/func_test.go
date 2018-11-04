// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/ast"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parser"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("Func", func() {

	It("should call the function and return the result", func() {
		fs := parsley.NewFileSet()
		expectedCtx := parsley.NewContext(
			fs,
			&parsleyfakes.FakeReader{},
		)
		expectedLeftRecCtx := data.NewIntMap(map[int]int{1: 2})
		expectedPos := parsley.Pos(2)
		expectedCurtailingParsers := data.NewIntSet(1)
		expectedNodes := ast.NewTerminalNode("x", nil, "", parsley.Pos(1), parsley.Pos(2))
		expectedErr := &parsleyfakes.FakeError{}
		expectedErr.ErrorReturns("test error")

		var actualCtx *parsley.Context
		var actualLeftRecCtx data.IntMap
		var actualPos parsley.Pos
		var actualErr parsley.Error
		p := parser.Func(func(ctx *parsley.Context, leftRecCtx data.IntMap, pos parsley.Pos) (parsley.Node, data.IntSet, parsley.Error) {
			actualCtx = ctx
			actualLeftRecCtx = leftRecCtx
			actualPos = pos
			return expectedNodes, expectedCurtailingParsers, expectedErr
		})

		actualNodes, actualCurtailingParsers, actualErr := p.Parse(expectedCtx, expectedLeftRecCtx, expectedPos)

		Expect(actualCtx).To(BeIdenticalTo(expectedCtx))
		Expect(actualCurtailingParsers).To(Equal(expectedCurtailingParsers))
		Expect(actualNodes).To(Equal(expectedNodes))
		Expect(actualLeftRecCtx).To(Equal(expectedLeftRecCtx))
		Expect(actualPos).To(BeIdenticalTo(expectedPos))
		Expect(actualErr).To(BeEquivalentTo(expectedErr))
	})
})
