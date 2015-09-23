package apperror_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/theduke/go-apperror"
)

var _ = Describe("Apperror", func() {
	It("Should .AddError()", func() {
		err1 := &AppError{
			Code:   "nested",
			Errors: []error{errors.New("nested1"), errors.New("nested2")},
		}

		err2 := &AppError{Code: "parent"}
		err2.AddError(err1)

		Expect(len(err2.Errors)).To(Equal(3))
		Expect(err2.Errors[0].(*AppError).Code).To(Equal("nested"))
		Expect(err2.Errors[1].Error()).To(Equal("nested1"))
		Expect(err2.Errors[2].Error()).To(Equal("nested2"))
	})

	It("Should .New() with message", func() {
		Expect(New("code", "msg").Message).To(Equal("msg"))
	})

	It("Should .New() with public flag", func() {
		Expect(New("code", true).Public).To(BeTrue())
	})

	It("Should .New() with data", func() {
		Expect(New("code", 444).Data).To(Equal(444))
	})

	It("Should .New() with errors", func() {
		Expect(New("code", []error{errors.New("nested")}).Errors[0]).To(Equal(errors.New("nested")))
	})

	It("Should .New() with message, flag, data and errors", func() {
		err := New("code", "msg", true, 11, []error{errors.New("nested")})
		Expect(err).To(Equal(&AppError{
			Code:    "code",
			Message: "msg",
			Public:  true,
			Data:    11,
			Errors:  []error{errors.New("nested")},
		}))
	})

	It("Should .WrapError with just a code", func() {
		err := WrapError(errors.New("nested"), "err")
		Expect(err).To(Equal(&AppError{
			Code:    "err",
			Message: "nested",
			Errors:  []error{errors.New("nested")},
		}))
	})

	It("Should .WrapError with code and message", func() {
		err := WrapError(errors.New("nested"), "err", "message")
		Expect(err).To(Equal(&AppError{
			Code:    "err",
			Message: "message: nested",
			Errors:  []error{errors.New("nested")},
		}))
	})

	It("Should .WrapError with code and public", func() {
		err := WrapError(errors.New("nested"), "err", true)
		Expect(err).To(Equal(&AppError{
			Code:   "err",
			Public: true,
			Errors: []error{errors.New("nested")},
		}))
	})

	It("Should .WrapError with code and data", func() {
		err := WrapError(errors.New("nested"), "err", 55)
		Expect(err).To(Equal(&AppError{
			Code:    "err",
			Message: "nested",
			Data:    55,
			Errors:  []error{errors.New("nested")},
		}))
	})

	It("Should .WrapError with code and message and public", func() {
		err := WrapError(errors.New("nested"), "err", "message", true)
		Expect(err).To(Equal(&AppError{
			Code:    "err",
			Message: "message",
			Public:  true,
			Errors:  []error{errors.New("nested")},
		}))
	})
})
