package state

import (
	"context"
	"encoding/json"
	"sync"
)

type ErrStackInterface interface {
	error
	HasErrs() bool
	GetErrs() []Error
}

type ListenerInterface interface {
	OnFinish(st OpInterface)
	OnError(st OpInterface)
	BeforeMove(ctx context.Context, st OpInterface, from []string, to []string) error
	AfterMove(ctx context.Context, st OpInterface, from []string, to []string)
}

// State of net.
// If you need custom serialization you can use this struct embedded in your implementation.
type State struct {
	places     map[string]struct{}
	errStack   *ErrStack
	isFinished bool
	listener   ListenerInterface
	mu         sync.Mutex
}

// NewState init new state.
func NewState() *State {
	return &State{
		places:     make(map[string]struct{}),
		errStack:   NewErrStack(),     // We can init inside value object without DI.
		listener:   NewStubListener(), // We can init inside value object without DI. And set it after if need.
		isFinished: false,
	}
}

// WithListener set listener to state.
func (s *State) WithListener(listener ListenerInterface) {
	s.listener = listener
}

// GetErrorStack returns errStack from state.
func (s *State) GetErrorStack() ErrStackInterface {
	return s.errStack
}

// GetPlaces returns places list copy.
func (s *State) GetPlaces() []string {
	res := make([]string, 0, len(s.places))
	for place := range s.places {
		res = append(res, place)
	}

	return res
}

// IsError return true if that is errStack state.
func (s *State) IsError() bool {
	return s.errStack.HasErrs()
}

// AddError state.
// If try to set nil errStack, panic will happen.
func (s *State) AddError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.errStack.Add(BuildError(err))

	s.listener.OnError(s)
}

// IsFinished the net.
func (s *State) IsFinished() bool {
	return s.isFinished
}

// IsFinished the net.
func (s *State) SetFinished() error {
	if s.IsFinished() {
		return NewError(ErrCodeStateIsAlreadyFinished, "Can't set finished state, because state is already finished")
	}

	s.isFinished = true

	s.listener.OnFinish(s)

	return nil
}

// IsStarted the net.
func (s *State) IsStarted() bool {
	return len(s.places) > 0
}

// MoveTokensFromPlacesToPlaces for create new state.
func (s *State) MoveTokensFromPlacesToPlaces(ctx context.Context, from []string, to []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.IsError() {
		return NewError(ErrCodeStateIsErrorState, "Can't process state to new places, state is errStack")
	}

	if s.IsFinished() {
		return NewError(ErrCodeStateIsFinished, "Can't process state to new places, state is finished")
	}

	for _, place := range from {
		if _, ok := s.places[place]; !ok {
			return NewErrorf(
				ErrCodeStateHasNotTokenInPlace,
				"State has not token in place '%s', state places: %+v",
				place, s.places,
			)
		}
	}

	for _, place := range to {
		if _, ok := s.places[place]; ok {
			return NewErrorf(
				ErrCodeStateAlreadyHasTokenInPlace,
				"State already has token in place '%s', state places: %+v",
				place, s.places,
			)
		}
	}

	if err := s.listener.BeforeMove(ctx, s, from, to); err != nil {
		return err
	}

	for _, place := range from {
		delete(s.places, place)
	}

	for _, place := range to {
		s.places[place] = struct{}{}
	}

	s.listener.AfterMove(ctx, s, from, to)

	return nil
}

type jsonState struct {
	Places     []string  `json:"places"`
	ErrStack   *ErrStack `json:"errStack"`
	IsFinished bool      `json:"isFinished"`
}

func (s *State) MarshalJSON() ([]byte, error) {
	jsonPlaces := make([]string, 0, len(s.places))

	for place := range s.places {
		jsonPlaces = append(jsonPlaces, place)
	}

	jsonSt := jsonState{
		Places:     jsonPlaces,
		ErrStack:   s.errStack,
		IsFinished: s.isFinished,
	}

	return json.Marshal(jsonSt)
}

func (s *State) UnmarshalJSON(data []byte) error {
	var jsonSt jsonState

	if err := json.Unmarshal(data, &jsonSt); err != nil {
		return err
	}

	s.places = make(map[string]struct{})

	for _, place := range jsonSt.Places {
		s.places[place] = struct{}{}
	}

	s.errStack = jsonSt.ErrStack
	s.isFinished = jsonSt.IsFinished

	if s.listener == nil {
		s.listener = NewStubListener()
	}

	return nil
}
