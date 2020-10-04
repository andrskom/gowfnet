package gowfnet

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStubListener_BeforeStart_AlwaysNil(t *testing.T) {
	assert.NoError(t, NewStubListener().BeforeStart(context.Background()))
}

func TestStubListener_BeforeTransition_AlwaysNil(t *testing.T) {
	ctrl := gomock.NewController(t)

	assert.NoError(t, NewStubListener().BeforeTransition(context.Background(), "", NewMockStateOpInterface(ctrl)))
}

func TestStubListener_HasStateListener_AlwaysFalse(t *testing.T) {
	assert.False(t, NewStubListener().HasStateListener())
}

func TestStubListener_GetStateListener_AlwaysNil(t *testing.T) {
	assert.Nil(t, NewStubListener().GetStateListener())
}
