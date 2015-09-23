package apperror

type Error interface {
	GetCode() string
	GetMessage() string
	GetData() interface{}

	IsPublic() bool

	GetErrors() []error
	SetErrors(errs []error)
	AddError(err error)

	Error() string
}
