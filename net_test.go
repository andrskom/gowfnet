package gowfnet

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
	mocks "github.com/andrskom/gowfnet/moscks"
	"github.com/andrskom/gowfnet/state"
)

func TestCtxSubject(t *testing.T) {
	eSubj := &struct{}{}

	ctx := SetSubject(context.Background(), eSubj)
	aSubj, ok := GetSubject(ctx)
	assert.True(t, ok)
	assert.Same(t, eSubj, aSubj)
}

func TestGetSubject_NotSetSubject_ReturnsNotOk(t *testing.T) {
	aSubj, ok := GetSubject(context.Background())
	assert.False(t, ok)
	assert.Nil(t, aSubj)
}

func TestNet_Start_AlreadyStarted_ReturnsErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := mocks.NewMockInterface(ctrl)
	config.EXPECT().GetTransitions().Return(cfg.MinimalTransitionRegistry{})
	config.EXPECT().GetPlaces().Return([]cfg.IDGetter{})

	net := NewNet(config)

	st := NewMockStateInterface(ctrl)
	st.EXPECT().IsStarted().Return(true)

	err := net.Start(context.Background(), st)
	assert.Error(t, err)
	tErr, ok := err.(*state.Error)
	assert.True(t, ok)
	assert.Equal(
		t,
		state.ErrCode("gowfnet.state.alreadyStarted"),
		tErr.GetCode(),
	)
	assert.Equal(
		t,
		"State already started in net",
		tErr.GetMessage(),
	)
}

func TestNet_Process_MoveToStartErr_ExpectedErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := mocks.NewMockInterface(ctrl)
	config.EXPECT().GetTransitions().Return(cfg.MinimalTransitionRegistry{})
	config.EXPECT().GetPlaces().Return([]cfg.IDGetter{})
	config.EXPECT().GetStart().Return(cfg.StringID("b"))

	net := NewNet(config)
	st := NewMockStateInterface(ctrl)
	st.EXPECT().IsStarted().Return(false)

	eErr := errors.New("a")
	st.EXPECT().MoveTokensFromPlacesToPlaces(context.Background(), []string{}, []string{"b"}).Return(eErr)

	err := net.Start(context.Background(), st)
	assert.Same(t, eErr, err)
}

func TestNet_Process_StartToFirstPlaceIsFinishWithoutErr_ReturnsNil(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := mocks.NewMockInterface(ctrl)
	config.EXPECT().GetTransitions().Return(cfg.MinimalTransitionRegistry{})
	config.EXPECT().GetPlaces().Return([]cfg.IDGetter{cfg.StringID("b")})
	config.EXPECT().GetStart().Return(cfg.StringID("b"))
	config.EXPECT().GetFinish().Return(cfg.StringID("b"))

	net := NewNet(config)
	st := NewMockStateInterface(ctrl)
	st.EXPECT().IsStarted().Return(false)
	st.EXPECT().MoveTokensFromPlacesToPlaces(context.Background(), []string{}, []string{"b"}).Return(nil)
	st.EXPECT().SetFinished().Return(nil)

	err := net.Start(context.Background(), st)
	assert.NoError(t, err)
}

func TestNet_Process_StartToFirstPlaceIsFinishWithErr_ReturnsTheSameErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := mocks.NewMockInterface(ctrl)
	config.EXPECT().GetTransitions().Return(cfg.MinimalTransitionRegistry{})
	config.EXPECT().GetPlaces().Return([]cfg.IDGetter{cfg.StringID("b")})
	config.EXPECT().GetStart().Return(cfg.StringID("b"))
	config.EXPECT().GetFinish().Return(cfg.StringID("b"))

	net := NewNet(config)
	st := NewMockStateInterface(ctrl)
	st.EXPECT().IsStarted().Return(false)
	st.EXPECT().MoveTokensFromPlacesToPlaces(context.Background(), []string{}, []string{"b"}).Return(nil)

	eErr := errors.New("a")
	st.EXPECT().SetFinished().Return(eErr)

	err := net.Start(context.Background(), st)
	assert.Same(t, eErr, err)
}

func TestNet_Process_StartToOnePlaceIsNotFinish_NoErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := mocks.NewMockInterface(ctrl)
	config.EXPECT().GetTransitions().Return(cfg.MinimalTransitionRegistry{})
	config.EXPECT().GetPlaces().Return([]cfg.IDGetter{cfg.StringID("b")})
	config.EXPECT().GetStart().Return(cfg.StringID("b"))
	config.EXPECT().GetFinish().Return(cfg.StringID("c"))

	net := NewNet(config)
	st := NewMockStateInterface(ctrl)
	st.EXPECT().IsStarted().Return(false)
	st.EXPECT().MoveTokensFromPlacesToPlaces(context.Background(), []string{}, []string{"b"}).Return(nil)

	err := net.Start(context.Background(), st)
	assert.NoError(t, err)
}

func TestNet_Process_StartToSomePlaces_NoErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := mocks.NewMockInterface(ctrl)
	config.EXPECT().GetTransitions().Return(cfg.MinimalTransitionRegistry{
		"t": {
			To:   []cfg.StringID{"e", "f"},
			From: []cfg.StringID{"h"},
		},
	})
	config.EXPECT().GetPlaces().Return([]cfg.IDGetter{cfg.StringID("b")})
	config.EXPECT().GetStart().Return(cfg.StringID("b"))
	config.EXPECT().GetFinish().Return(cfg.StringID("c"))

	net := NewNet(config)
	st := NewMockStateInterface(ctrl)
	st.EXPECT().IsStarted().Return(true)
	st.EXPECT().MoveTokensFromPlacesToPlaces(context.Background(), []string{"h"}, []string{"e", "f"}).Return(nil)

	err := net.Transit(context.Background(), st, "t")
	assert.NoError(t, err)
}

func TestNet_Transit_NotStarted_ExpectedErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := mocks.NewMockInterface(ctrl)
	config.EXPECT().GetTransitions().Return(cfg.MinimalTransitionRegistry{
		"t": {
			To:   []cfg.StringID{"e", "f"},
			From: []cfg.StringID{"h"},
		},
	})
	config.EXPECT().GetPlaces().Return([]cfg.IDGetter{cfg.StringID("b")})

	net := NewNet(config)
	st := NewMockStateInterface(ctrl)
	st.EXPECT().IsStarted().Return(false)

	err := net.Transit(context.Background(), st, "t")
	assert.Equal(t, state.NewError(state.ErrCodeStateIsNotStarted, "Can't transit, state is not started"), err)
}

func TestNet_Transit_UnknownTransition_ExpectedErr(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := mocks.NewMockInterface(ctrl)
	config.EXPECT().GetTransitions().Return(cfg.MinimalTransitionRegistry{
		"t": {
			To:   []cfg.StringID{"e", "f"},
			From: []cfg.StringID{"h"},
		},
	})
	config.EXPECT().GetPlaces().Return([]cfg.IDGetter{cfg.StringID("b")})

	net := NewNet(config)
	st := NewMockStateInterface(ctrl)
	st.EXPECT().IsStarted().Return(true)

	err := net.Transit(context.Background(), st, "t1")
	assert.Equal(
		t,
		state.NewErrorf(
			state.ErrCodeNetDoesntKnowAboutTransition,
			"Net doesn't know about transition '%s'",
			"t1",
		),
		err,
	)
}

func TestState_WithListener(t *testing.T) {
	ctrl := gomock.NewController(t)

	config := mocks.NewMockInterface(ctrl)
	config.EXPECT().GetTransitions().Return(cfg.MinimalTransitionRegistry{})
	config.EXPECT().GetPlaces().Return([]cfg.IDGetter{})

	net := NewNet(config)

	listener := NewStubListener()
	net.WithListener(listener)

	assert.Same(t, listener, net.listener)
}
