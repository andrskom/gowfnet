package validator

import (
	"fmt"
	"strings"
)

// Error is struct for validation errors.
type Error struct {
	list []string
}

// NewError init error.
func NewError() *Error {
	return &Error{
		list: make([]string, 0),
	}
}

// BuildErrorf return new created err with added msg.
func BuildErrorf(format string, args ...interface{}) *Error {
	err := NewError()
	err.Addf(format, args...)

	return err
}

// Addf message to err list.
func (e *Error) Addf(format string, args ...interface{}) {
	e.list = append(e.list, fmt.Sprintf(format, args...))
}

// Has errors.
func (e *Error) Has() bool {
	return len(e.list) > 0
}

// Has errors.
func (e *Error) Get() []string {
	return e.list
}

func (e *Error) Error() string {
	return " - " + strings.Join(e.list, "\n - ") + "\n"
}

func PrepareResultErr(err *Error) error {
	if err.Has() {
		return err
	}

	return nil
}
