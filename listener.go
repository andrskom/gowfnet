package gowfnet

import (
	"context"

	"github.com/andrskom/gowfnet/cfg"
)

type (
	AfterPlaceListenerFunc func(ctx context.Context, config cfg.Interface, state StateReadInterface, placeID string)

	BeforeTransitionListenerFunc func(
		ctx context.Context,
		config cfg.Interface,
		state StateReadInterface,
		transitionID string,
	) error
)

type ListenerOp func(l *Listener)

func BeforeTransitionOp(listenerFunc BeforeTransitionListenerFunc) ListenerOp {
	return func(l *Listener) {
		l.beforeTransition = append(l.beforeTransition, listenerFunc)
	}
}

func AfterPlaceOp(listenerFunc AfterPlaceListenerFunc) ListenerOp {
	return func(l *Listener) {
		l.afterPlace = append(l.afterPlace, listenerFunc)
	}
}

type Listener struct {
	afterPlace       []AfterPlaceListenerFunc
	beforeTransition []BeforeTransitionListenerFunc
}

func NewListener(ops ...ListenerOp) *Listener {
	l := Listener{
		afterPlace:       make([]AfterPlaceListenerFunc, 0),
		beforeTransition: make([]BeforeTransitionListenerFunc, 0),
	}

	for _, op := range ops {
		op(&l)
	}

	return &l
}

func (l *Listener) BeforeTransition(
	ctx context.Context,
	config cfg.Interface,
	state StateReadInterface,
	transitionID string,
) error {
	for _, fn := range l.beforeTransition {
		if err := fn(ctx, config, state, transitionID); err != nil {
			return err
		}
	}

	return nil
}

func (l *Listener) AfterPlace(ctx context.Context, config cfg.Interface, state StateReadInterface, placeID string) {
	for _, fn := range l.afterPlace {
		fn(ctx, config, state, placeID)
	}
}
