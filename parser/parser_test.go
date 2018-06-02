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
		expectedHistory := &parsleyfakes.FakeHistory{}
		expectedLeftRecCtx := data.NewIntMap(map[int]int{1: 2})
		expectedReader := &parsleyfakes.FakeReader{}
		expectedPos := parsley.Pos(2)
		expectedCurtailingParsers := data.NewIntSet(1)
		expectedNodes := ast.NewTerminalNode("x", nil, parsley.Pos(1), parsley.Pos(2))
		expectedErr := parsley.NewErrorf(parsley.Pos(1), "testerr")
		var actualHistory parsley.History
		var actualLeftRecCtx data.IntMap
		var actualReader parsley.Reader
		var actualPos parsley.Pos
		p := parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (parsley.Node, parsley.Error, data.IntSet) {
			actualHistory = h
			actualLeftRecCtx = leftRecCtx
			actualReader = r
			actualPos = pos
			return expectedNodes, expectedErr, expectedCurtailingParsers
		})

		actualNodes, actualErr, actualCurtailingParsers := p.Parse(expectedHistory, expectedLeftRecCtx, expectedReader, expectedPos)

		Expect(actualHistory).To(BeIdenticalTo(expectedHistory))
		Expect(actualCurtailingParsers).To(Equal(expectedCurtailingParsers))
		Expect(actualNodes).To(Equal(expectedNodes))
		Expect(actualLeftRecCtx).To(Equal(expectedLeftRecCtx))
		Expect(actualReader).To(BeIdenticalTo(expectedReader))
		Expect(actualErr).To(BeIdenticalTo(expectedErr))
	})

	Describe("WithName", func() {
		It("should create a named parser", func() {
			f := parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (parsley.Node, parsley.Error, data.IntSet) {
				return nil, nil, data.EmptyIntSet
			})
			p := f.WithName("p1")
			Expect(p.Name()).To(Equal("p1"))
		})

		Context("when a function is passed as name", func() {
			It("should call the function to get the name", func() {
				f := parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (parsley.Node, parsley.Error, data.IntSet) {
					return nil, nil, data.EmptyIntSet
				})
				name := func() string {
					return "p1"
				}
				p := f.WithName(name)
				Expect(p.Name()).To(Equal("p1"))
			})
		})

		Context("called with an invalid parameter", func() {
			It("should panic", func() {
				f := parser.Func(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos parsley.Pos) (parsley.Node, parsley.Error, data.IntSet) {
					return nil, nil, data.EmptyIntSet
				})
				Expect(func() { f.WithName(nil) }).To(Panic())
			})
		})
	})
})
