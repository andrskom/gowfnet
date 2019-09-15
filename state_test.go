package gowfnet

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewState(t *testing.T) {
	require.Equal(
		t,
		&State{
			places:     make(map[string]struct{}),
			error:      nil,
			isFinished: false,
			mu:         sync.Mutex{},
		},
		NewState(),
	)
}

func TestState_SetError(t *testing.T) {
	require.Equal(
		t,
		&State{
			places:     make(map[string]struct{}),
			error:      nil,
			isFinished: false,
			mu:         sync.Mutex{},
		},
		NewState(),
	)
}
