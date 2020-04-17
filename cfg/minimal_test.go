package cfg

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andrskom/gowfnet"
)

func TestCreateStringID(t *testing.T) {
	assert.Equal(t, StringID("a"), CreateStringID("a"))
}

func TestStringID_GetID(t *testing.T) {
	assert.Equal(t, "a", CreateStringID("a").GetID())
}

func TestMinimal_GetStart(t *testing.T) {
	cfg := Minimal{Start: "a"}
	assert.Equal(t, CreateStringID("a"), cfg.GetStart())
}

func TestMinimal_GetFinish(t *testing.T) {
	cfg := Minimal{Finish: "a"}
	assert.Equal(t, CreateStringID("a"), cfg.GetFinish())
}

func TestMinimal_GetPlaces(t *testing.T) {
	cfg := Minimal{
		Places: []StringID{
			"a",
		},
	}
	assert.Equal(t,
		[]IDGetter{
			CreateStringID("a"),
		},
		cfg.GetPlaces(),
	)
}

func TestMinimal_GetTransitions(t *testing.T) {
	cfg := Minimal{
		Transitions: MinimalTransitionRegistry{
			"a": {
				To:   []StringID{"b"},
				From: nil,
			},
		},
	}
	assert.Equal(t,
		map[string]TransitionInterface{
			"a": MinimalTransition{
				To:   []StringID{"b"},
				From: nil,
			},
		},
		cfg.GetTransitions().GetAsMap(),
	)
}

func TestMinimalTransitionRegistry_GetByID_TransitionNotSet_ExpectedErr(t *testing.T) {
	reg := MinimalTransitionRegistry{}

	res, err := reg.GetByID(StringID("a"))

	assert.Nil(t, res)
	gowfnet.AssertErrCodeEqual(t, ErrCodeUnknownTransitionID, err)
}

func TestMinimalTransitionRegistry_GetByID_TransitionSet_ExpectedVal(t *testing.T) {
	reg := MinimalTransitionRegistry{
		"a": {
			To:   []StringID{"b"},
			From: nil,
		},
	}

	res, err := reg.GetByID(StringID("a"))

	assert.NoError(t, err)
	assert.Equal(t,
		MinimalTransition{
			To:   []StringID{"b"},
			From: nil,
		},
		res,
	)
}

func TestMinimalTransition_GetFrom(t *testing.T) {
	tr := MinimalTransition{
		To:   nil,
		From: []StringID{"a"},
	}

	assert.Equal(t,
		[]IDGetter{CreateStringID("a")},
		tr.GetFrom(),
	)
}

func TestMinimalTransition_GetTo(t *testing.T) {
	tr := MinimalTransition{
		To:   []StringID{"a"},
		From: nil,
	}

	assert.Equal(t,
		[]IDGetter{CreateStringID("a")},
		tr.GetTo(),
	)
}

func TestMinimalTransitionRegistry_Marshalling(t *testing.T) {
	reg := MinimalTransitionRegistry{
		"a": {
			To:   []StringID{"b"},
			From: nil,
		},
	}

	bytes, err := json.Marshal(reg)
	require.NoError(t, err)

	var res MinimalTransitionRegistry

	err = json.Unmarshal(bytes, &res)
	assert.NoError(t, err)
	assert.Equal(t, reg, res)
}
