package state

import "context"

type OpInterface interface {
	GetErrorStack() ErrStackInterface
	GetPlaces() []string
	IsError() bool
	AddError(err error)
	IsFinished() bool
	SetFinished() error
	IsStarted() bool
}

type StubListener struct {
}

func NewStubListener() *StubListener {
	return &StubListener{}
}

func (l *StubListener) OnFinish(st OpInterface) {}

func (l *StubListener) OnError(st OpInterface) {}

func (l *StubListener) BeforeMove(ctx context.Context, st OpInterface, from []string, to []string) error {
	return nil
}

func (l *StubListener) AfterMove(ctx context.Context, st OpInterface, from []string, to []string) {}
