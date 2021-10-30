// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parsley_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/parsley"
)

var _ = Describe("NilPosition", func() {
	It("implements the Positon interface", func() {
		var _ parsley.Position = parsley.NilPosition
	})

	It("returns with a non-empty string representation", func() {
		Expect(parsley.NilPosition.String()).ToNot(BeEmpty())
	})
})
