package gowfnet

import (
	"encoding/json"
	"fmt"
	"testing"
)

type ErrCode string

const (
	ErrCodeStateHasNotTokenInPlace      = "gowfnet.stateHasNotTokenInPlace"     // nolint:gosec
	ErrCodeStateAlreadyHasTokenInPlace  = "gowfnet.stateAlreadyHasTokenInPlace" // nolint:gosec
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

// ErrStack is a stack of errors for state.
// We use that instead simple error because sometimes we need to register many errors in resp.
type ErrStack struct {
	stack []Error
}

// NewErrStack init err stack.
func NewErrStack() *ErrStack {
	return &ErrStack{stack: make([]Error, 0)}
}

// Add err to stack.
// If you send nil *Error, panic will happen.
func (s *ErrStack) Add(err *Error) {
	if err == nil {
		panic("must use not nil err")
	}

	s.stack = append(s.stack, *err)
}

// HasErrs in stack.
func (s *ErrStack) HasErrs() bool {
	return len(s.stack) > 0
}

// GetErrs fro stack.
func (s *ErrStack) GetErrs() []Error {
	return s.stack
}

// Error interface implementation.
func (s *ErrStack) Error() string {
	res := ""

	for i := 0; i < len(s.stack); i++ {
		res += fmt.Sprintf("%d) %s;\n", i, s.stack[i].Error())
	}

	return res
}

type jsonErrStack struct {
	Stack []Error `json:"stack"`
}

func (s *ErrStack) UnmarshalJSON(data []byte) error {
	var res jsonErrStack

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	s.stack = res.Stack

	return nil
}

func (s ErrStack) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonErrStack{Stack: s.stack})
}

// Error is err model of component.
type Error struct {
	code    ErrCode
	message string
}

// NewError init errStack.
func NewError(errorCode ErrCode, msg string) *Error {
	return &Error{
		code:    errorCode,
		message: msg,
	}
}

// NewErrorf init errStack by format.
func NewErrorf(errorCode ErrCode, format string, args ...interface{}) *Error {
	return NewError(errorCode, fmt.Sprintf(format, args...))
}

// ErrorIs check errStack.
func (e *Error) Is(errorCode ErrCode) bool {
	return e.code == errorCode
}

func (e *Error) Error() string {
	return e.message
}

// ErrorIs compare err and errStack code. If err is not nil, is *Error type and have the same code or false.
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

// BuildError from errStack interface.
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

type jsonErr struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func (e Error) MarshalJSON() ([]byte, error) {
	jsonSt := jsonErr{
		Code:    e.code,
		Message: e.message,
	}

	return json.Marshal(jsonSt)
}

func (e *Error) UnmarshalJSON(data []byte) error {
	var jsonErr jsonErr

	if err := json.Unmarshal(data, &jsonErr); err != nil {
		return err
	}

	e.code = jsonErr.Code
	e.message = jsonErr.Message

	return nil
}

func AssertErrCodeEqual(t *testing.T, code ErrCode, err error) {
	if ErrorIs(code, err) {
		return
	}

	t.Errorf(`ecode of err is not equal
expected: %s
actualErr: %#v`, code, err)
	t.Fail()
}

func RequireErrCodeEqual(t *testing.T, code ErrCode, err error) {
	if ErrorIs(code, err) {
		return
	}

	t.Errorf(`ecode of err is not equal
expected: %s
actualErr: %#v`, code, err)
	t.FailNow()
}
