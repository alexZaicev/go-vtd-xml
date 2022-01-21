package erroring

import "fmt"

// InternalErrorType for usage with std library errors.As
var InternalErrorType = &InternalError{}

// InvalidArgumentErrorType for usage with std library errors.As
var InvalidArgumentErrorType = &InvalidArgumentError{}

// ParseErrorType for usage with std library errors.As
var ParseErrorType = &ParseError{}

// baseError represents a generic error from the domain package that provides
// 'error' functionality to the rest of the typed errors in the package.
type baseError struct {
	Err error
	msg string
}

// newBaseError is a constructor for a base error. It should not be called
// directly outside constructing other errors in this package.
func newBaseError(msg string, err error) baseError {
	return baseError{
		Err: err,
		msg: msg,
	}
}

// Error allows baseError and any structs that embed it to satisfy the error
// interface.
func (e *baseError) Error() string {
	return e.msg
}

// Unwrap allows baseError and any structs that embed it to be used with the
// error wrapping utilities introduced in go 1.13.
func (e *baseError) Unwrap() error {
	// This nil check accounts for the situation where the embedded *baseError
	// in one of the public errors is nil - if it has been constructed without
	// using one of the helper functions (e.g in other package's unit tests).
	if e == nil {
		return nil
	}
	return e.Err
}

// InvalidArgumentError represents a parameter that does not meet expected criteria.
// It is differentiated from InvalidInput in that it is expected to be used for
// defensive programming (e.g. nil checks) rather than for validating external
// input into the system.
type InvalidArgumentError struct {
	baseError
	Parameter string
	Msg       string
}

// NewInvalidArgumentError constructs a new InvalidArgumentError error struct. An error
// can be specified to wrap, but is not expected in most cases.
func NewInvalidArgumentError(param, msg string, err error) *InvalidArgumentError {
	return &InvalidArgumentError{
		baseError: newBaseError(
			fmt.Sprintf("invalid argument %s: %s", param, msg),
			err,
		),
		Parameter: param,
		Msg:       msg,
	}
}

// InternalError encompasses logic errors that are not immediately resolvable
// and that the caller is not expected to perform any actions beyond returning
// an error itself with a generic error message. For example, the scanning of a
// SQL row that fails would result in an internal error as either there is a bug
// in the code or the data has been corrupted.
type InternalError struct {
	baseError
}

// NewInternalError constructs a new InternalError, wrapping the provided error.
func NewInternalError(msg string, err error) *InternalError {
	return &InternalError{
		baseError: newBaseError(
			fmt.Sprintf("an internal error occurred: %s", msg),
			err,
		),
	}
}

// ParseError represent a some sort od parsing error that may occur trying to
// convert non-numeric string to a number.
type ParseError struct {
	baseError
	Msg string
}

// NewParseError constructs a new ParseError, wrapping the provided error.
func NewParseError(msg string, err error) *InternalError {
	return &InternalError{
		baseError: newBaseError(
			fmt.Sprintf("a parse error occurred: %s", msg),
			err,
		),
	}
}
