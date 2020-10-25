// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("NewError", func() {

	var (
		err    parsley.Error
		cause  error
		format string
		values []interface{}
		pos    parsley.Pos
	)

	BeforeEach(func() {
		format = "some %s"
		values = []interface{}{"error"}
		pos = parsley.Pos(1)
		cause = fmt.Errorf(format, values...)
	})

	JustBeforeEach(func() {
		err = parsley.NewError(pos, cause)
	})

	It("implements error", func() {
		var _ error = err
	})

	Describe("Pos()", func() {
		It("returns with the position", func() {
			Expect(err.Pos()).To(BeIdenticalTo(pos))
		})
	})

	Describe("Error()", func() {
		It("returns with the formatted error message", func() {
			Expect(err.Error()).To(Equal("some error"))
		})
	})

	Describe("Cause()", func() {
		It("returns with the original error", func() {
			Expect(err.Cause()).To(BeIdenticalTo(cause))
		})
	})

	Context("When the cause is a parlsey error", func() {
		BeforeEach(func() {
			cause = parsley.NewErrorf(parsley.Pos(2), "parsley error")
		})

		It("should return the cause and not wrap it", func() {
			Expect(err).To(BeIdenticalTo(cause))
		})
	})
})

var _ = Describe("NewErrorf", func() {

	var (
		err    parsley.Error
		format string
		values []interface{}
		pos    parsley.Pos
	)

	BeforeEach(func() {
		format = "some %s"
		values = []interface{}{"error"}
		pos = parsley.Pos(1)
	})

	JustBeforeEach(func() {
		err = parsley.NewErrorf(pos, format, values...)
	})

	It("implements error", func() {
		var _ error = err
	})

	Describe("Pos()", func() {
		It("returns with the position", func() {
			Expect(err.Pos()).To(BeIdenticalTo(pos))
		})
	})

	Describe("Error()", func() {
		It("returns with the formatted error message", func() {
			Expect(err.Error()).To(Equal("some error"))
		})
	})

	Describe("Cause()", func() {
		It("returns with the original error", func() {
			Expect(err.Cause()).To(MatchError(errors.New("some error")))
		})
	})
})
