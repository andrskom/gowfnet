package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
)

func TestNewDuplicatedPlacesInTransitions(t *testing.T) {
	assert.Equal(t, &DuplicatedPlacesInTransitions{}, NewDuplicatedPlacesInTransitions())
}

func TestDuplicatedPlacesInTransitions_Validate_DuplicatedInFrom_ExpectedErr(t *testing.T) {
	minCfg := cfg.Minimal{
		Transitions: map[string]cfg.MinimalTransition{
			"a": {
				From: []cfg.StringID{"b", "b", "c"},
			},
		},
	}

	v := NewDuplicatedPlacesInTransitions()

	assert.Equal(
		t,
		BuildErrorf("place with id 'b' is duplicated in transition with id 'a' in section from"),
		v.Validate(minCfg),
	)
}

func TestDuplicatedPlacesInTransitions_Validate_DuplicatedInTo_ExpectedErr(t *testing.T) {
	minCfg := cfg.Minimal{
		Transitions: map[string]cfg.MinimalTransition{
			"a": {
				To: []cfg.StringID{"b", "b", "c"},
			},
		},
	}

	v := NewDuplicatedPlacesInTransitions()

	assert.Equal(
		t,
		BuildErrorf("place with id 'b' is duplicated in transition with id 'a' in section to"),
		v.Validate(minCfg),
	)
}

func TestDuplicatedPlacesInTransitions_Validate_CorrectCfg_NoErr(t *testing.T) {
	minCfg := cfg.Minimal{
		Transitions: map[string]cfg.MinimalTransition{
			"a": {
				To:   []cfg.StringID{"b"},
				From: []cfg.StringID{"c"},
			},
		},
	}

	v := NewDuplicatedPlacesInTransitions()

	assert.NoError(t, v.Validate(minCfg))
}
