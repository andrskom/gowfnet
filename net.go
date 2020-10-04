package gowfnet

import (
	"context"

	"github.com/andrskom/gowfnet/cfg"
	"github.com/andrskom/gowfnet/state"
)

// nolint:gochecknoglobals
var ctxSubject = &struct{}{}

func SetSubject(ctx context.Context, subj interface{}) context.Context {
	return context.WithValue(ctx, ctxSubject, subj)
}

func GetSubject(ctx context.Context) (interface{}, bool) {
	data := ctx.Value(ctxSubject)
	if data == nil {
		return nil, false
	}

	return data, true
}

type StateReadInterface interface {
	IsStarted() bool
	IsFinished() bool
	IsError() bool
	GetErrorStack() state.ErrStackInterface
	GetPlaces() []string
}

type StateOpInterface interface {
	StateReadInterface
	SetFinished() error
	AddError(err error)
}

type StateInterface interface {
	StateOpInterface
	WithListener(listener state.ListenerInterface)
	MoveTokensFromPlacesToPlaces(ctx context.Context, from []string, to []string) error
}

type ListenerInterface interface {
	BeforeStart(ctx context.Context) error
	AfterStart(ctx context.Context)
	BeforeTransition(ctx context.Context, transitionID string, state StateOpInterface) error
	AfterTransition(ctx context.Context, transitionID string, state StateOpInterface)
	HasStateListener() bool
	GetStateListener() state.ListenerInterface
}

type Net struct {
	cfg           cfg.Interface
	placeMap      map[string]cfg.IDGetter
	transitionMap map[string]cfg.TransitionInterface
	listener      ListenerInterface
}

func NewNet(config cfg.Interface) *Net {
	net := &Net{
		cfg:           config,
		placeMap:      make(map[string]cfg.IDGetter),
		transitionMap: config.GetTransitions().GetAsMap(),
		listener:      NewStubListener(),
	}

	for _, place := range config.GetPlaces() {
		net.placeMap[place.GetID()] = place
	}

	return net
}

func (n *Net) WithListener(listener ListenerInterface) {
	n.listener = listener
}

// Start workflow net.
//
// Use ctx for cancel operation and send subject of operation.
func (n *Net) Start(ctx context.Context, s StateInterface) error {
	if s.IsStarted() {
		return state.NewError(state.ErrCodeStateAlreadyStarted, "State already started in net")
	}

	if err := n.listener.BeforeStart(ctx); err != nil {
		return err
	}

	defer n.listener.AfterStart(ctx)

	return n.process(ctx, s, []string{}, buildStringSliceFromIDGetter(n.cfg.GetStart()))
}

// Transit to new places(state).
//
// Use ctx for cancel operation and send subject of operation.
func (n *Net) Transit(ctx context.Context, s StateInterface, transitionID string) error {
	if !s.IsStarted() {
		return state.NewError(state.ErrCodeStateIsNotStarted, "Can't transit, state is not started")
	}

	transition, ok := n.transitionMap[transitionID]
	if !ok {
		return state.NewErrorf(
			state.ErrCodeNetDoesntKnowAboutTransition,
			"Net doesn't know about transition '%s'",
			transitionID,
		)
	}

	if err := n.listener.BeforeTransition(ctx, transitionID, s); err != nil {
		return err
	}

	err := n.process(
		ctx,
		s,
		buildStringSliceFromIDGetter(transition.GetFrom()...),
		buildStringSliceFromIDGetter(transition.GetTo()...),
	)

	if err != nil {
		return err
	}

	n.listener.AfterTransition(ctx, transitionID, s)

	return nil
}

// nolint:unparam
func (n *Net) process(ctx context.Context, s StateInterface, fromPlaces []string, toPlaces []string) error {
	if n.listener.HasStateListener() {
		s.WithListener(n.listener.GetStateListener())
	}

	if err := s.MoveTokensFromPlacesToPlaces(ctx, fromPlaces, toPlaces); err != nil {
		return err
	}

	if len(toPlaces) != 1 {
		return nil
	}

	if toPlaces[0] == n.cfg.GetFinish().GetID() {
		return s.SetFinished()
	}

	return nil
}

func buildStringSliceFromIDGetter(in ...cfg.IDGetter) []string {
	res := make([]string, 0, len(in))
	for _, id := range in {
		res = append(res, id.GetID())
	}

	return res
}
