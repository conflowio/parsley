// Copyright (c) 2017 Opsidian Ltd.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package terminal_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/conflowio/parsley/data"
	"github.com/conflowio/parsley/parsley"
	"github.com/conflowio/parsley/text"
	"github.com/conflowio/parsley/text/terminal"
)

var _ = Describe("TimeDuration", func() {

	var p = terminal.TimeDuration("time.Duration")

	DescribeTable("should match",
		func(input string, startPos int, value time.Duration, nodePos int, endPos int) {
			f := text.NewFile("textfile", []byte(input))
			fs := parsley.NewFileSet(f)
			r := text.NewReader(f)
			ctx := parsley.NewContext(fs, r)
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, parsley.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Token()).To(Equal("TIME_DURATION"))
			Expect(res.Schema()).To(Equal("time.Duration"))
			Expect(res.(parsley.LiteralNode).Value()).To(Equal(value))
			Expect(res.Pos()).To(Equal(parsley.Pos(nodePos)))
			Expect(res.ReaderPos()).To(Equal(parsley.Pos(endPos)))
		},
		Entry("5s beginning", "5s ---", 1, 5*time.Second, 1, 3),
		Entry("5s middle", "--- 5s ---", 5, 5*time.Second, 5, 7),
		Entry("5s end", "--- 5s", 5, 5*time.Second, 5, 7),
		Entry("0s", "0s", 1, time.Duration(0), 1, 3),
		Entry("5ns", "5ns", 1, 5*time.Nanosecond, 1, 4),
		Entry("5us", "5us", 1, 5*time.Microsecond, 1, 4),
		Entry("5µs, U+00B5 = micro symbol", "5µs", 1, 5*time.Microsecond, 1, 5),
		Entry("5μs, U+03BC = Greek letter mu", "5μs", 1, 5*time.Microsecond, 1, 5),
		Entry("5ms", "5ms", 1, 5*time.Millisecond, 1, 4),
		Entry("5s", "5s", 1, 5*time.Second, 1, 3),
		Entry("5m", "5m", 1, 5*time.Minute, 1, 3),
		Entry("5h", "5h", 1, 5*time.Hour, 1, 3),
		Entry("+5s", "+5s", 1, 5*time.Second, 1, 4),
		Entry("-5s", "-5s", 1, -5*time.Second, 1, 4),
		Entry("1.5m", "1.5m", 1, 90*time.Second, 1, 5),
		Entry("1h30m45s", "1h30m45s", 1, time.Hour+30*time.Minute+45*time.Second, 1, 9),
	)

	DescribeTable("should not match",
		func(input string, startPos int) {
			f := text.NewFile("textfile", []byte(input))
			fs := parsley.NewFileSet(f)
			r := text.NewReader(f)
			ctx := parsley.NewContext(fs, r)
			res, curtailingParsers, err := p.Parse(ctx, data.EmptyIntMap, f.Pos(startPos))
			Expect(curtailingParsers).To(Equal(data.EmptyIntSet))
			Expect(err).To(HaveOccurred())
			Expect(err.Cause()).To(MatchError("was expecting time duration"))
			Expect(err.Pos()).To(Equal(f.Pos(startPos)))
			Expect(res).To(BeNil())
		},
		Entry("empty", "", 0),
		Entry("a", "a", 0),
		Entry("-", "-", 0),
		Entry("+", "+", 0),
		Entry("integer", "0", 0),
		Entry("float", "0.1", 0),
		Entry("float without fractional", "1.s", 0),
	)
})
