package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
)

type validatorMock struct {
	err      error
	callsNum int
}

func (v *validatorMock) Validate(c cfg.Interface) error {
	v.callsNum++

	return v.err
}

func TestNew_NoValidators_EmptyList(t *testing.T) {
	com := New()
	assert.Equal(t,
		&Component{validators: make([]Validator, 0)},
		com,
	)
}

func TestNew_WithValidators_ExpectedComponent(t *testing.T) {
	v := &validatorMock{err: errors.New("a")}
	com := New(v)
	assert.Equal(t,
		&Component{validators: []Validator{v}},
		com,
	)
}

func TestComponent_Validate_ErrNotHappened_NoErr(t *testing.T) {
	v := &validatorMock{}
	com := New(v)
	assert.NoError(t, com.Validate(&cfg.Minimal{}))
	assert.Equal(t, 1, v.callsNum, "unexpected numbers of mock calls")
}

func TestComponent_Validate_ErrHappened_Err(t *testing.T) {
	v := &validatorMock{err: errors.New("a")}
	com := New(v)
	assert.Equal(t, errors.New("a"), com.Validate(&cfg.Minimal{}))
	assert.Equal(t, 1, v.callsNum, "unexpected numbers of mock calls")
}
