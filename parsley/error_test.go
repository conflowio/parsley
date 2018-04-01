package parsley_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/parsley"
)

var _ = Describe("Error", func() {

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
		err = parsley.NewError(pos, format, values...)
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
})

var _ = Describe("WrapError", func() {

	var (
		err        parsley.Error
		wrappedErr parsley.Error
		pos        parsley.Pos
		format     string
		values     []interface{}
	)

	BeforeEach(func() {
		pos = parsley.Pos(1)
		wrappedErr = parsley.NewError(pos, "some error")
		format = "I wrap {{err}} as a %s"
		values = []interface{}{"test"}
	})

	JustBeforeEach(func() {
		err = parsley.WrapError(wrappedErr, format, values...)
	})

	It("should return with an error with the given position", func() {
		Expect(err.Pos()).To(BeIdenticalTo(pos))
	})

	It("should replace {{err}} and the placeholders in the error message", func() {
		Expect(err.Error()).To(Equal("I wrap some error as a test"))
	})

	Context("when the error is not wrapped in a message", func() {
		BeforeEach(func() {
			format = ""
			values = []interface{}{}
		})
		It("should create the same error as NewError would", func() {
			Expect(err).To(Equal(wrappedErr))
		})
	})
})
