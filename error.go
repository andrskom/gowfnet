package gowfnet

import (
	"fmt"
)

type ErrCode string

// nolint:gosec
const (
	ErrCodeStateHasNotTokenInPlace      = "gowfnet.stateHasNotTokenInPlace"
	ErrCodeStateAlreadyHasTokenInPlace  = "gowfnet.stateAlreadyHasTokenInPlace"
	ErrCodeNetDoesntKnowAboutTransition = "gowfnet.netDoesntKnowAboutTransition"
	ErrCodeNetDoesntKnowAboutPlace      = "gowfnet.netDoesntKnowAboutPlace"
	ErrCodeUnknown                      = "gowfnet.unknown"
	ErrCodeStateAlreadyStarted          = "gowfnet.stateAlreadyStarted"
	ErrCodeStateIsNotStarted            = "gowfnet.stateIsNotStarted"
	ErrCodeStateIsFinished              = "gowfnet.stateIsFinished"
	ErrCodeStateIsErrorState            = "gowfnet.stateIsErrorState"
	ErrCodeRegistryNetAlreadyRegistered = "gowfnet.registryNetAlreadyRegistered"
	ErrCodeRegistryNetNotRegistered     = "gowfnet.registryNetNotRegistered"
)

// Error is err model of component.
type Error struct {
	code    ErrCode
	message string
}

// NewError init error.
func NewError(errorCode ErrCode, msg string) *Error {
	return &Error{
		code:    errorCode,
		message: msg,
	}
}

// NewErrorf init error by format.
func NewErrorf(errorCode ErrCode, format string, args ...interface{}) *Error {
	return NewError(errorCode, fmt.Sprintf(format, args...))
}

// ErrorIs check error.
func (e *Error) Is(errorCode ErrCode) bool {
	return e.code == errorCode
}

func (e *Error) Error() string {
	return e.message
}

// ErrorIs compare err and error code. If err is not nil, is *Error type and have the same code or false.
func ErrorIs(code ErrCode, err error) bool {
	if err == nil {
		return false
	}
	errModel, ok := err.(*Error)
	if !ok {
		return false
	}
	return errModel.Is(code)
}

// BuildError from error interface.
// If arg contains nil, return nil.
// If arg contains Error type, return that.
// If arg contains another type, build Error with unknown code.
func BuildError(err error) *Error {
	if err == nil {
		return nil
	}
	res, ok := err.(*Error)
	if ok {
		return res
	}

	return NewError(ErrCodeUnknown, err.Error())
}
