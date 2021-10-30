// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parser"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/parsley/parsleyfakes"
)

var _ = Describe("ReturnError", func() {

	var q *parsleyfakes.FakeParser
	var ctx *parsley.Context
	var leftRecCtx data.IntMap
	var pos parsley.Pos
	var res, qres parsley.Node
	var err, qerr parsley.Error
	var cp, qcp data.IntSet

	JustBeforeEach(func() {
		q = &parsleyfakes.FakeParser{}
		q.ParseReturnsOnCall(0, qres, qcp, qerr)
		p := parser.ReturnError(q, errors.New("custom error"))
		res, cp, err = p.Parse(ctx, leftRecCtx, pos)
	})

	BeforeEach(func() {
		ctx = parsley.NewContext(parsley.NewFileSet(), &parsleyfakes.FakeReader{})
		leftRecCtx = data.NewIntMap(map[int]int{1: 2})
		pos = parsley.Pos(1)
		qres = nil
		qcp = data.NewIntSet(1, 2)
		qerr = nil
	})

	Context("when q returns no error", func() {

		BeforeEach(func() {
			qres = &parsleyfakes.FakeNode{}
		})

		It("should return the original result", func() {
			Expect(res).To(BeEquivalentTo(qres))
			Expect(cp).To(BeEquivalentTo(qcp))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when q returns a not found error with the reader's position", func() {

		BeforeEach(func() {
			qerr = parsley.NewError(pos, parsley.NotFoundError("some error"))
		})

		It("should return the custom error", func() {
			Expect(res).To(BeNil())
			Expect(cp).To(BeEquivalentTo(qcp))
			Expect(err).To(MatchError("custom error"))
		})
	})

	Context("when q returns a parse error with the reader's position", func() {

		BeforeEach(func() {
			qerr = parsley.NewErrorf(pos, "some error")
		})

		It("should return the original error", func() {
			Expect(res).To(BeNil())
			Expect(cp).To(BeEquivalentTo(qcp))
			Expect(err).To(MatchError("some error"))
		})
	})

	Context("when q returns an error with higher position than the reader's position", func() {

		BeforeEach(func() {
			qerr = parsley.NewErrorf(parsley.Pos(2), "some error")
		})

		It("should return the original error", func() {
			Expect(res).To(BeNil())
			Expect(cp).To(BeEquivalentTo(qcp))
			Expect(err).To(MatchError("some error"))
		})
	})

	Context("when q returns no result and no error", func() {
		It("should return the custom error", func() {
			Expect(res).To(BeNil())
			Expect(cp).To(BeEquivalentTo(qcp))
			Expect(err).To(MatchError("custom error"))
		})
	})

})
