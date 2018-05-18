// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package text_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/text"
)

var _ = Describe("Position", func() {
	var (
		pos          *text.Position
		line, column int
		filename     string
	)

	BeforeEach(func() {
		line = 1
		column = 2
		filename = "testfile"
	})

	JustBeforeEach(func() {
		pos = text.NewPosition(filename, line, column)
	})

	It("should implement the parsley.Position interface", func() {
		var _ parsley.Position = text.Position{}
	})

	It("should return with a string containing all information", func() {
		Expect(pos.String()).To(Equal("testfile:1:2"))
	})

	Context("no filename", func() {
		BeforeEach(func() {
			filename = ""
		})

		It("should return with a string without the filename", func() {
			Expect(pos.String()).To(Equal("1:2"))
		})
	})
})
