package parsley_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opsidian/parsley/parsley"
	"github.com/opsidian/parsley/parsley/parsleyfakes"
)

var _ = Describe("Error", func() {

	var (
		err   *parsley.Error
		cause error
		pos   parsley.Position
	)

	BeforeEach(func() {
		cause = errors.New("some error")
		pos = &parsleyfakes.FakePosition{}
		pos.(*parsleyfakes.FakePosition).StringReturns("testpos")
	})

	JustBeforeEach(func() {
		err = parsley.NewError(cause, pos)
	})

	It("implements error", func() {
		var _ error = &parsley.Error{}
	})

	Describe("NewError()", func() {
		Context("when created with the same error type", func() {
			BeforeEach(func() {
				cause = parsley.NewError(errors.New("some error"), &parsleyfakes.FakePosition{})
			})

			It("should return the original error instead creating a new one", func() {
				Expect(err).To(Equal(cause))
			})
		})
	})

	Describe("Cause()", func() {
		It("returns with the cause", func() {
			Expect(err.Cause()).To(Equal(cause))
		})
	})

	Describe("Pos()", func() {
		It("returns with the position", func() {
			Expect(err.Pos()).To(BeIdenticalTo(pos))
		})
	})

	Describe("Error()", func() {
		It("returns with a formatted error message", func() {
			Expect(err.Error()).To(Equal("some error at testpos"))
		})

		Context("when created with NilPosition", func() {
			BeforeEach(func() {
				pos = parsley.NilPosition
			})

			It("returns with the error message without a position", func() {
				Expect(err.Error()).To(Equal("some error"))
			})
		})

		Context("when created with nil position", func() {
			BeforeEach(func() {
				pos = nil
			})

			It("returns with the error message without a position", func() {
				Expect(err.Error()).To(Equal("some error"))
			})
		})
	})
})

var _ = Describe("WrapError", func() {

	var (
		err    *parsley.Error
		cause  *parsley.Error
		pos    parsley.Position
		format string
		values []interface{}
	)

	BeforeEach(func() {
		pos = &parsleyfakes.FakePosition{}
		pos.(*parsleyfakes.FakePosition).StringReturns("testpos")
		cause = parsley.NewError(errors.New("some error"), pos)
		format = "I wrap {{err}} as a %s"
		values = []interface{}{"test"}
		pos.(*parsleyfakes.FakePosition).StringReturns("testpos")
	})

	JustBeforeEach(func() {
		err = parsley.WrapError(cause, format, values...)
	})

	It("should return with an error with the given position", func() {
		Expect(err.Pos()).To(BeIdenticalTo(pos))
	})

	It("should replace {{err}} and the placeholders in the error message", func() {
		Expect(err.Cause()).To(MatchError("I wrap some error as a test"))
	})

	It("should return with a wrapped error message and position", func() {
		Expect(err.Error()).To(Equal("I wrap some error as a test at testpos"))
	})

	Context("when the error is not wrapped in a message", func() {
		BeforeEach(func() {
			format = ""
			values = []interface{}{}
		})
		It("should create the same error as NewError would", func() {
			Expect(err).To(Equal(cause))
		})
	})
})
