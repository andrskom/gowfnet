package gowfnet

import (
	"context"

	"github.com/andrskom/gowfnet/state"
)

type StubListener struct{}

func NewStubListener() *StubListener {
	return &StubListener{}
}

func (l *StubListener) BeforeStart(ctx context.Context) error {
	return nil
}

func (l *StubListener) AfterStart(ctx context.Context) {}

func (l *StubListener) BeforeTransition(ctx context.Context, transitionID string, state StateOpInterface) error {
	return nil
}

func (l *StubListener) AfterTransition(ctx context.Context, transitionID string, state StateOpInterface) {
}

func (l *StubListener) HasStateListener() bool {
	return false
}

func (l *StubListener) GetStateListener() state.ListenerInterface {
	return nil
}
