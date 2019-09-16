package gowfnet

import (
	"encoding/json"
	"sync"
)

// State of net.
type State struct {
	places     map[string]struct{}
	error      *Error
	isFinished bool
	mu         sync.Mutex
}

// NewState init new state.
func NewState() *State {
	return &State{
		places:     make(map[string]struct{}),
		error:      nil,
		isFinished: false,
	}
}

// GetError return error from state.
func (s *State) GetError() *Error {
	return s.error
}

// GetPlaces return list of places.
func (s *State) GetPlaces() []string {
	res := make([]string, 0, len(s.places))
	for place := range s.places {
		res = append(res, place)
	}

	return res
}

// IsError return true if that is error state.
func (s *State) IsError() bool {
	return s.error != nil
}

// SetError state.
// If try to set nil error, panic will happen.
// If try to set error state while state is finished, panic will happen.
// If try to set error state while state is not started, panic will happen.
// If is already error state, not set new error state.
func (s *State) SetError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err == nil {
		panic("arguments of function can't be nil")
	}

	if !s.IsStarted() {
		panic("state is not started")
	}

	if s.IsFinished() {
		panic("state already is finished")
	}

	if !s.IsError() {
		s.error = BuildError(err)
	}
}

// IsFinished the net.
func (s *State) IsFinished() bool {
	return s.isFinished
}

// IsStarted the net.
func (s *State) IsStarted() bool {
	return len(s.places) > 0
}

// HasTokensInPlaces check token in places.
func (s *State) HasTokensInPlaces(places []string) bool {
	for _, place := range places {
		if _, ok := s.places[place]; !ok {
			return false
		}
	}

	return true
}

// MoveTokensFromPlacesToPlaces for create new state.
func (s *State) MoveTokensFromPlacesToPlaces(from []string, to []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.IsError() {
		return NewError(ErrCodeStateIsErrorState, "Can't process state to new places, state is error")
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

	for _, place := range from {
		delete(s.places, place)
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

	for _, place := range to {
		s.places[place] = struct{}{}
	}

	return nil
}

type jsonState struct {
	Places     []string `json:"places"`
	Error      *Error   `json:"error"`
	IsFinished bool     `json:"isFinished"`
}

// nolint:govet
func (s State) MarshalJSON() ([]byte, error) {
	jsonPlaces := make([]string, 0, len(s.places))
	for place := range s.places {
		jsonPlaces = append(jsonPlaces, place)
	}
	jsonSt := jsonState{
		Places:     jsonPlaces,
		Error:      s.error,
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
	s.error = jsonSt.Error
	s.isFinished = jsonSt.IsFinished
	return nil
}
