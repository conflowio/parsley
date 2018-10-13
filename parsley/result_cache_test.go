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

var _ = Describe("Result Cache", func() {
	var (
		rc parsley.ResultCache
	)

	BeforeEach(func() {
		rc = parsley.NewResultCache()
	})

	Describe("SaveResult/GetResult", func() {
		var (
			parserIndex           int
			pos                   parsley.Pos
			res, res1, res2, res3 *parsley.Result
			found                 bool
			leftRecCtx            data.IntMap
		)

		BeforeEach(func() {
			res1 = &parsley.Result{
				CurtailingParsers: data.EmptyIntSet,
				Error:             parsley.NewErrorf(parsley.Pos(2), "some error"),
				LeftRecCtx:        data.EmptyIntMap,
				Node:              &parsleyfakes.FakeNode{},
			}
			res2 = &parsley.Result{
				CurtailingParsers: data.EmptyIntSet,
				Error:             nil,
				LeftRecCtx:        data.EmptyIntMap,
				Node:              &parsleyfakes.FakeNode{},
			}
			res3 = &parsley.Result{
				CurtailingParsers: data.EmptyIntSet,
				Error:             nil,
				LeftRecCtx:        data.EmptyIntMap,
				Node:              &parsleyfakes.FakeNode{},
			}

			leftRecCtx = data.EmptyIntMap
			parserIndex = 1
			pos = parsley.Pos(1)
		})

		JustBeforeEach(func() {
			rc.Save(1, parsley.Pos(1), res1)
			rc.Save(1, parsley.Pos(2), res2)
			rc.Save(2, parsley.Pos(1), res3)
			res, found = rc.Get(parserIndex, pos, leftRecCtx)
		})

		It("should return previously saved result", func() {
			Expect(res).To(BeEquivalentTo(res1))
			Expect(found).To(BeTrue())
		})

		Context("when getting a result for the same parser but different position", func() {
			BeforeEach(func() {
				pos = parsley.Pos(2)
			})
			It("should return a different result", func() {
				Expect(res).To(BeEquivalentTo(res2))
				Expect(found).To(BeTrue())
			})
		})

		Context("when getting a result for a different parser but the same position", func() {
			BeforeEach(func() {
				parserIndex = 2
			})
			It("should return a different result", func() {
				Expect(res).To(BeEquivalentTo(res3))
				Expect(found).To(BeTrue())
			})
		})

		Context("when there is no result for the parser", func() {
			BeforeEach(func() {
				parserIndex = 99
			})
			It("should return nil", func() {
				Expect(res).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when there is no result for the position", func() {
			BeforeEach(func() {
				pos = parsley.Pos(99)
			})
			It("should return nil", func() {
				Expect(res).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when there is no result for the position", func() {
			BeforeEach(func() {
				pos = parsley.Pos(99)
			})
			It("should return nil", func() {
				Expect(res).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})

		Context("when in the saved result any of the left-rec counts are higher for a parser", func() {
			BeforeEach(func() {
				res1.LeftRecCtx = data.NewIntMap(map[int]int{1: 2, 2: 1})
				leftRecCtx = data.NewIntMap(map[int]int{1: 1, 2: 1})
			})
			It("should not be returned", func() {
				Expect(res).To(BeNil())
				Expect(found).To(BeFalse())
			})
		})
	})
})
