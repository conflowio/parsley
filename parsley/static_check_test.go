// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("StaticCheck", func() {
	var (
		r        *parsleyfakes.FakeReader
		ctx      *parsley.Context
		node     *parsleyfakes.FakeStaticCheckableNode
		evalCtx  interface{}
		checkErr parsley.Error
		err      error
	)

	BeforeEach(func() {
		f := &parsleyfakes.FakeFile{}
		position := &parsleyfakes.FakePosition{}
		position.StringReturns("testpos")
		f.PositionReturns(position)

		fs := parsley.NewFileSet(f)

		r = &parsleyfakes.FakeReader{}
		ctx = parsley.NewContext(fs, r)
		r.PosReturns(parsley.Pos(1))
		node = &parsleyfakes.FakeStaticCheckableNode{}
		evalCtx = "context"
		checkErr = nil
		err = nil
	})

	JustBeforeEach(func() {
		node.StaticCheckReturns(checkErr)
		err = parsley.StaticCheck(ctx, node, evalCtx)
	})

	Context("if the static check succeeds", func() {
		It("should return no error", func() {
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("if the static check has an error", func() {
		BeforeEach(func() {
			err := &parsleyfakes.FakeError{}
			err.PosReturns(parsley.Pos(1))
			err.ErrorReturns("some error")
			checkErr = err
		})
		It("should return an error", func() {
			Expect(err).To(MatchError("some error at testpos"))
		})
	})
})
