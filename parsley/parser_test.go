// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/data"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("ParserFunc", func() {

	It("should call the function and return the result", func() {
		expectedHistory := &parsleyfakes.FakeHistory{}
		expectedLeftRecCtx := data.NewIntMap(map[int]int{1: 2})
		expectedReader := &parsleyfakes.FakeReader{}
		expectedPos := 2
		expectedCurtailingParsers := data.NewIntSet(1)
		expectedNodes := []parsley.Node{nil}
		expectedErr := parsley.NewError(parsley.Pos(1), "testerr")
		var actualHistory parsley.History
		var actualLeftRecCtx data.IntMap
		var actualReader parsley.Reader
		var actualPos int
		p := parsley.ParserFunc(func(h parsley.History, leftRecCtx data.IntMap, r parsley.Reader, pos int) (data.IntSet, []parsley.Node, parsley.Error) {
			actualHistory = h
			actualLeftRecCtx = leftRecCtx
			actualReader = r
			actualPos = pos
			return expectedCurtailingParsers, expectedNodes, expectedErr
		})

		actualCurtailingParsers, actualNodes, actualErr := p.Parse(expectedHistory, expectedLeftRecCtx, expectedReader, expectedPos)

		Expect(actualHistory).To(BeIdenticalTo(expectedHistory))
		Expect(actualCurtailingParsers).To(Equal(expectedCurtailingParsers))
		Expect(actualNodes).To(Equal(expectedNodes))
		Expect(actualLeftRecCtx).To(Equal(expectedLeftRecCtx))
		Expect(actualReader).To(BeIdenticalTo(expectedReader))
		Expect(actualErr).To(BeIdenticalTo(expectedErr))
	})

})
