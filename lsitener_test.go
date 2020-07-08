package gowfnet

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/andrskom/gowfnet/cfg"
	"github.com/andrskom/gowfnet/state"
)

func TestBeforeTransitionOp(t *testing.T) {
	l := Listener{
		afterPlace:       nil,
		beforeTransition: make([]BeforeTransitionListenerFunc, 0),
	}

	eErr := errors.New("")
	fn := func(
		ctx context.Context,
		config cfg.Interface,
		state StateReadInterface,
		transitionID string,
	) error {
		return eErr
	}

	BeforeTransitionOp(fn)(&l)

	assert.Len(t, l.beforeTransition, 1)
	assert.Same(t, eErr, l.beforeTransition[0](context.Background(), &cfg.Minimal{}, &state.State{}, ""))
}

func TestAfterPlaceOp(t *testing.T) {
	l := Listener{
		afterPlace:       nil,
		beforeTransition: make([]BeforeTransitionListenerFunc, 0),
	}

	check := 0
	fn := func(
		ctx context.Context,
		config cfg.Interface,
		state StateReadInterface,
		placeID string,
	) {
		check++
	}

	AfterPlaceOp(fn)(&l)

	assert.Len(t, l.afterPlace, 1)
	l.afterPlace[0](context.Background(), &cfg.Minimal{}, &state.State{}, "")
	assert.Equal(t, 1, check, "check the listener func through a call failed")
}

type MockListeners struct {
	mock.Mock
}

func (m *MockListeners) AfterPlace(
	ctx context.Context,
	config cfg.Interface,
	state StateReadInterface,
	placeID string,
) {
	m.Called(ctx, config, state, placeID)
}

func (m *MockListeners) BeforeTransition(
	ctx context.Context,
	config cfg.Interface,
	state StateReadInterface,
	transitionID string,
) error {
	args := m.Called(ctx, config, state, transitionID)
	return args.Error(0)
}

func TestListener_BeforeTransition_ErrInFirst_ExpectedErr(t *testing.T) {
	m1 := &MockListeners{}
	m2 := &MockListeners{}

	ctx := context.Background()
	config := &cfg.Minimal{}
	state := &state.State{}
	transitionID := "a"

	eErr := errors.New("")
	m1.On("BeforeTransition", ctx, config, state, transitionID).Return(eErr).Once()

	listener := NewListener(BeforeTransitionOp(m1.BeforeTransition), BeforeTransitionOp(m2.BeforeTransition))
	err := listener.BeforeTransition(ctx, config, state, transitionID)

	assert.Same(t, eErr, err)
	m1.AssertNumberOfCalls(t, "BeforeTransition", 1)
	m2.AssertNotCalled(t, "BeforeTransition", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestListener_BeforeTransition_NoErr_NoErr(t *testing.T) {
	m1 := &MockListeners{}
	m2 := &MockListeners{}

	ctx := context.Background()
	config := &cfg.Minimal{}
	state := &state.State{}
	transitionID := "a"

	m1.On("BeforeTransition", ctx, config, state, transitionID).Return(nil).Once()
	m2.On("BeforeTransition", ctx, config, state, transitionID).Return(nil).Once()

	listener := NewListener(BeforeTransitionOp(m1.BeforeTransition), BeforeTransitionOp(m2.BeforeTransition))
	err := listener.BeforeTransition(ctx, config, state, transitionID)

	assert.NoError(t, err)
	m1.AssertNumberOfCalls(t, "BeforeTransition", 1)
	m2.AssertNumberOfCalls(t, "BeforeTransition", 1)
}

func TestListener_AfterPlace_SetTwo_CalledTwo(t *testing.T) {
	m1 := &MockListeners{}
	m2 := &MockListeners{}

	ctx := context.Background()
	config := &cfg.Minimal{}
	state := &state.State{}
	placeID := "a"

	m1.On("AfterPlace", ctx, config, state, placeID).Return().Once()
	m2.On("AfterPlace", ctx, config, state, placeID).Return().Once()

	listener := NewListener(AfterPlaceOp(m1.AfterPlace), AfterPlaceOp(m2.AfterPlace))
	listener.AfterPlace(ctx, config, state, placeID)

	m1.AssertNumberOfCalls(t, "AfterPlace", 1)
	m2.AssertNumberOfCalls(t, "AfterPlace", 1)
}
