package structs

import (
	"fmt"

	"errors"
)

const (
	errPermissionDenied = "Permission denied"
	errInvalidInput     = "Invalid input"
	errNotFound         = "Resource not found"
	errInternal         = "Internal error"
)

var (
	// ErrPermissionDenied :
	ErrPermissionDenied = errors.New(errPermissionDenied)

	// ErrInvalidInput ...
	ErrInvalidInput = errors.New(errInvalidInput)

	// ErrInternal ...
	ErrInternal = errors.New(errInternal)

	// ErrNotFound ...
	ErrNotFound = errors.New(errNotFound)
)

// Error :
type Error struct {
	Message string
}

// NewError ...
func NewError(base error, extra ...interface{}) error {
	msg := base.Error()
	for _, v := range extra {
		msg = fmt.Sprintf("%s : %v", msg, v)
	}
	return &Error{
		Message: msg,
	}
}

func (e Error) Error() string {
	return e.Message
}

// NewInternalError :
func NewInternalError(msg string) error {
	return NewError(ErrInternal, msg)
}

// NewInvalidInputError :
func NewInvalidInputError(msg string) error {
	return NewError(ErrInvalidInput, msg)
}
