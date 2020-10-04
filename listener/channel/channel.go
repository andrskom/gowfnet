package channel

import (
	"context"
	"fmt"

	"github.com/andrskom/gowfnet"
	"github.com/andrskom/gowfnet/state"
)

type State struct {
	eventChan chan string
}

func NewState(eventChan chan string) *State {
	return &State{eventChan: eventChan}
}

func (s *State) OnFinish(st state.OpInterface) {
	s.eventChan <- "finished"
}

func (s *State) OnError(st state.OpInterface) {
	s.eventChan <- "errorHappened"
}

func (s *State) BeforeMove(ctx context.Context, st state.OpInterface, from []string, to []string) error {
	s.eventChan <- "move_" + s.printFromTo(from, to)

	return nil
}

func (s *State) AfterMove(ctx context.Context, st state.OpInterface, from []string, to []string) {
	s.eventChan <- "moved_" + s.printFromTo(from, to)
}

func (s *State) printFromTo(from []string, to []string) string {
	return fmt.Sprintf("FROM:%+v_TO:%+v", from, to)
}

type Listener struct {
	eventChan chan string
}

func New(chanSize int) *Listener {
	return &Listener{eventChan: make(chan string, chanSize)}
}

func (l *Listener) BeforeStart(ctx context.Context) error {
	l.eventChan <- "start"

	return nil
}

func (l *Listener) AfterStart(ctx context.Context) {
	l.eventChan <- "started"
}

func (l *Listener) BeforeTransition(ctx context.Context, transitionID string, state gowfnet.StateOpInterface) error {
	l.eventChan <- "transit_" + transitionID

	return nil
}

func (l *Listener) AfterTransition(ctx context.Context, transitionID string, state gowfnet.StateOpInterface) {
	l.eventChan <- transitionID + "_transited"
}

func (l *Listener) HasStateListener() bool {
	return true
}

func (l *Listener) GetStateListener() state.ListenerInterface {
	return NewState(l.eventChan)
}

func (l *Listener) ReadEvt() string {
	return <-l.eventChan
}
