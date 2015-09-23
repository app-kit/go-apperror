# go-apperror
Go error implementation for applications that supports nested errors.


Go s error  implementation is rather simplistic,  and leaves a lot to be desired
for applications that require more complex errors.

This package provides an `Error` interface and and an `AppError` implementation to 
provide rich errors, both internally and to a possibly public frontend.

The Error interface of course implements the go Error interface,
so you can always use use app errors as plain errors when required.

Errors can contain the following data:

* Code: a string code that uniquely identifies the error.
* Status: an integer value to identify the type of error (for example an http status code, or a unique error number).
* Message: an explicative message intended for humans.
* Data: arbitrary data related to the error.
* Errors: nested errors
* Public: specifies whether the error may be presented publically (for example to an api user or website)


## Usage

```go
import(
	"fmt"
	"github.com/theduke/go-apperror"
)

// Creating a new error.
func createString() (string, apperror.Error) {
	// Return an error with 'code' "my_application_error", 'status' 22 and a message.
	return "", apperror.New("my_application_error", 22, "Could not create string")
}

// Checking for errors.

func CreateString() string {
	s, err := createString()

	// Check if the error has a specific code.
	if apperror.IsCode("my_application_error") {
		fmt.Printf("XX error occured")
	} 

	// Check for error status.
	if apperror.IsStatus(33) {
		fmt.Printf("Error status 33")
  }

  return s
}

// Wrapping nested errors.

func dbQuery(q string) ([]string, apperror.Error) {
	result, err := db.Query(q)
	if err != nil {
		// Wrap the database error with a generic PUBLIC error.
		// The "true" argument makes the error public, which allows the error to be presented
		// to users.
		publicErr := apperror.Wrap(err, "db_query_error", "Database query failed", true)

		// The original error can be accessed with err.GetErrors()
		for _, err := range  publicErr.GetErrors() {
			fmt.Printf("Nested error: %v\n", err)
		}

		// Additional nested errors can be added.
		publicErr.AddError(errors.New("other_error"))
	}
	return result, nil
}

// Working with errors.

func handleError(err apperror.Error) {
	code := err.GetCode()
	status := err.GetStatus()
	msg := err.GetMessage()
	nestedErrors := err.GetErrors()
	data := err.GetData()

	isPublic := err.IsPublic()

	// Handle specific errors...
}
```

## Reference

* apperror.New()

Create a new error.

```go
// For apperror.New(), only the  code argument is  required.
// All additional arguments can be in any order.

err := apperror.New(code string, [status int,] [message string,] [data interface{}, ] [isPublic bool,] [nestedErrors []error])
```

* apperror.Wrap()

Wrap an error.

```go
// For apperror.Wrap(), the original error and the code argument are required.
// All additional arguments can be in any order.

 apperror.Wrap(originalError error, code string, [status int,] [message string,] [data interface{}, ] [isPublic bool,]) apperror.Error
```

* apperror.IsCode

Check if an error has a code.

```go
apperror.IsCode(err error, code string) bool
```

* apperror.IsStatus

Check if an error has a status.

```go
apperror.IsStatus(err error, status int) bool
```

* Manually creating errors.

You can also manually create the error.

```go
err := &apperror.AppError{
	Code: "code",
	Status: 111,
	Message: "msg",
	Data: []string{"1", "2"},
	Errors: []error{errors.New("xxx")},
	Public: true,
}
return err
```

## Marshal

AppError implements the json Marshal interface to provide more advanced marshalling for 
public errors:

The marshaled error will not include nested errors.

If marshaling a non-public error, all data will be stripped and 
the error will just contain the Code "app_error" and a generic error message.

## Tipps

* Always use the apperror.Error interface in function signatures, not the implementation apperror.AppError. This allows you to use different implementations if appropriate.
* Set errors intended for end users to Public, and check to only show public errors.
