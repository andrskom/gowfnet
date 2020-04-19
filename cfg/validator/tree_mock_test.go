package validator

import (
	"github.com/stretchr/testify/mock"

	"github.com/andrskom/gowfnet/cfg"
)

type BuilderMock struct {
	mock.Mock
}

func NewBuilderMock() *BuilderMock {
	return &BuilderMock{}
}

func (b *BuilderMock) Build(c cfg.Interface) (*Tree, error) {
	args := b.Called(c)
	return args.Get(0).(*Tree), args.Error(1)
}
