package apperror

import (
	"encoding/json"
	"fmt"
)

type AppError struct {
	Code    string      `json:"code,omitempty"`
	Status  int         `json:"status,omitempty"`
	Message string      `json:"title,omitempty"`
	Data    interface{} `json:"data"`
	Errors  []error     `json:"errors,omitempty"`
	Public  bool        `json:"-"`
}

// Ensure error implements the error interface.
var _ Error = (*AppError)(nil)

func (e AppError) GetCode() string {
	return e.Code
}

func (e AppError) GetStatus() int {
	return e.Status
}

func (e AppError) GetMessage() string {
	return e.Message
}

func (e AppError) GetData() interface{} {
	return e.Data
}

func (e AppError) IsPublic() bool {
	return e.Public
}

func (e AppError) GetErrors() []error {
	return e.Errors
}

func (e *AppError) SetErrors(errs []error) {
	for _, err := range errs {
		e.AddError(err)
	}
}

func (e *AppError) AddError(err error) {
	if appError, ok := err.(Error); ok {
		// If the error implements the Error interface,
		// merge the nested errors.
		errs := appError.GetErrors()
		appError.SetErrors(nil)

		e.Errors = append(e.Errors, appError)
		fmt.Printf("errs: %v\n", errs)
		if len(errs) > 0 {
			e.Errors = append(e.Errors, errs...)
		}
	} else {
		e.Errors = append(e.Errors, err)
	}
}

func (e AppError) Error() string {
	s := e.Code
	if e.Status != 0 {
		s += fmt.Sprintf("(%v)", e.Status)
	}
	if e.Message != "" {
		s += ": " + e.Message
	}

	if e.Data != nil {
		s += "\nError data: " + fmt.Sprintf("%+v", e.Data)
	}

	return s
}

func (e AppError) ToJson() []byte {
	js, err := e.MarshalJSON()
	if err != nil {
		return []byte(`{"code": "error_marshal_failed", "message": "Could not convert the returned error to json."}`)
	}

	return js
}

// Implement the json Marshaler interface.
func (e AppError) MarshalJSON() ([]byte, error) {
	e.Errors = nil
	if !e.Public {
		e.Code = "app_error"
		e.Status = 0
		e.Message = "An internal application error occurred"
		e.Data = nil
	}

	return json.Marshal(e)
}

// Create a new error. only required argument is string.
// Other arguments may be: a string to set the message, a bool to set
// the error to public, an int to set the status,
// a slice of errors to set the nested errors,
// and an arbitrary interface{} to set the error.Data.
func New(code string, args ...interface{}) *AppError {
	err := &AppError{
		Code: code,
	}

	for _, arg := range args {
		if status, ok := arg.(int); ok {
			err.Status = status
		} else if str, ok := arg.(string); ok {
			err.Message = str
		} else if flag, ok := arg.(bool); ok {
			err.Public = flag
		} else if errs, ok := arg.([]error); ok {
			err.SetErrors(errs)
		} else {
			err.Data = arg
		}
	}

	return err
}

// Wrap an error with an AppError.
// The required arguments are the error to wrap an an error code.
// Additionally you can supply another string argument as the message,
// a bool to set if the error is public, and an arbitrary interface{} value
// to set as data.
// If you do not supply a message, the original error will be converted to string
// and used as the message.
func Wrap(err error, code string, args ...interface{}) *AppError {
	wrap := New(code, args...)
	wrap.SetErrors(nil)
	wrap.AddError(err)

	if !wrap.Public {
		// For private errors, merge the custom message (if any) and the error message.
		msg := wrap.Message
		if msg != "" {
			msg += ": "
		}
		wrap.Message = msg + err.Error()
	}

	return wrap
}

func IsCode(err error, code string) bool {
	if appErr, ok := err.(Error); ok {
		return appErr.GetCode() == code
	}
	return err.Error() == code
}

func IsStatus(err error, status int) bool {
	if appErr, ok := err.(Error); ok {
		return appErr.GetStatus() == status
	}
	return false
}
