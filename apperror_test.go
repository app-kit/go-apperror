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

	It("Should .Error() with just a code", func() {
		Expect(AppError{Code: "code"}.Error()).To(Equal("code"))
	})

	It("Should .Error() with code and status", func() {
		Expect(AppError{Code: "code", Status: 100}.Error()).To(Equal("code(100)"))
	})

	It("Should .Error() with code, status and message", func() {
		Expect(AppError{Code: "code", Status: 100, Message: "msg"}.Error()).To(Equal("code(100): msg"))
	})

	It("Should .New() with status", func() {
		Expect(New("code", 100).Status).To(Equal(100))
	})

	It("Should .New() with message", func() {
		Expect(New("code", "msg").Message).To(Equal("msg"))
	})

	It("Should .New() with public flag", func() {
		Expect(New("code", true).Public).To(BeTrue())
	})

	It("Should .New() with data", func() {
		Expect(New("code", []string{}).Data).To(Equal([]string{}))
	})

	It("Should .New() with errors", func() {
		Expect(New("code", []error{errors.New("nested")}).Errors[0]).To(Equal(errors.New("nested")))
	})

	It("Should .New() with status, message, flag, data and errors", func() {
		err := New("code", 100, "msg", true, []string{}, []error{errors.New("nested")})
		Expect(err).To(Equal(&AppError{
			Code:    "code",
			Status:  100,
			Message: "msg",
			Public:  true,
			Data:    []string{},
			Errors:  []error{errors.New("nested")},
		}))
	})

	It("Should .WrapError with just a code", func() {
		err := Wrap(errors.New("nested"), "err")
		Expect(err).To(Equal(&AppError{
			Code:    "err",
			Message: "nested",
			Errors:  []error{errors.New("nested")},
		}))
	})

	It("Should .WrapError with code and message", func() {
		err := Wrap(errors.New("nested"), "err", "message")
		Expect(err).To(Equal(&AppError{
			Code:    "err",
			Message: "message: nested",
			Errors:  []error{errors.New("nested")},
		}))
	})

	It("Should .WrapError with code and public", func() {
		err := Wrap(errors.New("nested"), "err", true)
		Expect(err).To(Equal(&AppError{
			Code:   "err",
			Public: true,
			Errors: []error{errors.New("nested")},
		}))
	})

	It("Should .WrapError with code and data", func() {
		err := Wrap(errors.New("nested"), "err", []string{})
		Expect(err).To(Equal(&AppError{
			Code:    "err",
			Message: "nested",
			Data:    []string{},
			Errors:  []error{errors.New("nested")},
		}))
	})

	It("Should .WrapError with code, status, data, message and public", func() {
		err := Wrap(errors.New("nested"), "err", 100, "message", true, []string{})
		Expect(err).To(Equal(&AppError{
			Code:    "err",
			Status:  100,
			Message: "message",
			Data:    []string{},
			Public:  true,
			Errors:  []error{errors.New("nested")},
		}))
	})

	It("Should .IsCode() with plain error", func() {
		Expect(IsCode(errors.New("err"), "err")).To(BeTrue())
	})

	It("Should .IsCode() with AppError", func() {
		Expect(IsCode(New("err"), "err")).To(BeTrue())
	})

	It("Should .IsStatus() with AppError", func() {
		Expect(IsStatus(New("err", 100), 100)).To(BeTrue())
	})
})
