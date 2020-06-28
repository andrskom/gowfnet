package gowfnet

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
	"github.com/andrskom/gowfnet/mocks"
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
	config := &mocks.CfgInterface{}
	config.On("GetTransitions").Return(cfg.MinimalTransitionRegistry{})
	config.On("GetPlaces").Return([]cfg.IDGetter{})

	net := NewNet(config)

	st := &MockStateInterface{}
	st.On("IsStarted").Return(true).Once()

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
	config := &mocks.CfgInterface{}
	config.On("GetTransitions").Return(cfg.MinimalTransitionRegistry{})
	config.On("GetPlaces").Return([]cfg.IDGetter{})
	config.On("GetStart").Return(cfg.StringID("b"))

	net := NewNet(config)
	st := &MockStateInterface{}
	st.On("IsStarted").Return(false).Once()

	eErr := errors.New("a")
	st.On("MoveTokensFromPlacesToPlaces", []string{}, []string{"b"}).Return(eErr).Once()

	err := net.Start(context.Background(), st)
	assert.Same(t, eErr, err)
}

func TestNet_Process_StartToFirstPlaceIsFinishWithoutErr_ReturnsNil(t *testing.T) {
	config := &mocks.CfgInterface{}
	config.On("GetTransitions").Return(cfg.MinimalTransitionRegistry{})
	config.On("GetPlaces").Return([]cfg.IDGetter{cfg.StringID("b")})
	config.On("GetStart").Return(cfg.StringID("b"))
	config.On("GetFinish").Return(cfg.StringID("b"))

	net := NewNet(config)
	st := &MockStateInterface{}
	st.On("IsStarted").Return(false).Once()
	st.On("MoveTokensFromPlacesToPlaces", []string{}, []string{"b"}).Return(nil).Once()
	st.On("SetFinished").Return(nil).Once()

	err := net.Start(context.Background(), st)
	assert.NoError(t, err)
}

func TestNet_Process_StartToFirstPlaceIsFinishWithErr_ReturnsTheSameErr(t *testing.T) {
	config := &mocks.CfgInterface{}
	config.On("GetTransitions").Return(cfg.MinimalTransitionRegistry{})
	config.On("GetPlaces").Return([]cfg.IDGetter{cfg.StringID("b")})
	config.On("GetStart").Return(cfg.StringID("b"))
	config.On("GetFinish").Return(cfg.StringID("b"))

	net := NewNet(config)
	st := &MockStateInterface{}
	st.On("IsStarted").Return(false).Once()
	st.On("MoveTokensFromPlacesToPlaces", []string{}, []string{"b"}).Return(nil).Once()

	eErr := errors.New("a")
	st.On("SetFinished").Return(eErr).Once()

	err := net.Start(context.Background(), st)
	assert.Same(t, eErr, err)
}

func TestNet_Process_StartToOnePlaceIsNotFinish_NoErr(t *testing.T) {
	config := &mocks.CfgInterface{}
	config.On("GetTransitions").Return(cfg.MinimalTransitionRegistry{})
	config.On("GetPlaces").Return([]cfg.IDGetter{cfg.StringID("b")})
	config.On("GetStart").Return(cfg.StringID("b"))
	config.On("GetFinish").Return(cfg.StringID("c"))

	net := NewNet(config)
	st := &MockStateInterface{}
	st.On("IsStarted").Return(false).Once()
	st.On("MoveTokensFromPlacesToPlaces", []string{}, []string{"b"}).Return(nil).Once()

	err := net.Start(context.Background(), st)
	assert.NoError(t, err)
}

func TestNet_Process_StartToSomePlaces_NoErr(t *testing.T) {
	config := &mocks.CfgInterface{}
	config.On("GetTransitions").Return(cfg.MinimalTransitionRegistry{
		"t": {
			To:   []cfg.StringID{"e", "f"},
			From: []cfg.StringID{"h"},
		},
	})
	config.On("GetPlaces").Return([]cfg.IDGetter{cfg.StringID("b")})
	config.On("GetStart").Return(cfg.StringID("b"))
	config.On("GetFinish").Return(cfg.StringID("c"))

	net := NewNet(config)
	st := &MockStateInterface{}
	st.On("IsStarted").Return(true).Once()
	st.On("MoveTokensFromPlacesToPlaces", []string{"h"}, []string{"e", "f"}).Return(nil).Once()

	err := net.Transit(context.Background(), st, "t")
	assert.NoError(t, err)
}

func TestNet_Transit_NotStarted_ExpectedErr(t *testing.T) {
	config := &mocks.CfgInterface{}
	config.On("GetTransitions").Return(cfg.MinimalTransitionRegistry{
		"t": {
			To:   []cfg.StringID{"e", "f"},
			From: []cfg.StringID{"h"},
		},
	})
	config.On("GetPlaces").Return([]cfg.IDGetter{cfg.StringID("b")})

	net := NewNet(config)
	st := &MockStateInterface{}
	st.On("IsStarted").Return(false).Once()

	err := net.Transit(context.Background(), st, "t")
	assert.Equal(t, state.NewError(state.ErrCodeStateIsNotStarted, "Can't transit, state is not started"), err)
}

func TestNet_Transit_UnknownTransition_ExpectedErr(t *testing.T) {
	config := &mocks.CfgInterface{}
	config.On("GetTransitions").Return(cfg.MinimalTransitionRegistry{
		"t": {
			To:   []cfg.StringID{"e", "f"},
			From: []cfg.StringID{"h"},
		},
	})
	config.On("GetPlaces").Return([]cfg.IDGetter{cfg.StringID("b")})

	net := NewNet(config)
	st := &MockStateInterface{}
	st.On("IsStarted").Return(true).Once()

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
