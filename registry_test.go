package gowfnet

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ListenerMock struct {
	mock.Mock
}

func (m *ListenerMock) AfterPlaced(net *Net, state *State, placeID string, subject interface{}) {
	m.Called(net, subject, placeID, subject)
}

func TestWithGlobalListener(t *testing.T) {
	t.Run("not_nil_listener", func(t *testing.T) {
		listenerMock := &ListenerMock{}
		optsFunc := WithGlobalListener(listenerMock)
		opts := &RegistryBuildOpts{}
		optsFunc(opts)
		require.Equal(t, listenerMock, opts.GlobalListener)
	})
	t.Run("nil_listener", func(t *testing.T) {
		optsFunc := WithGlobalListener(nil)
		opts := &RegistryBuildOpts{}
		optsFunc(opts)
		require.Nil(t, opts.GlobalListener)
	})
}
