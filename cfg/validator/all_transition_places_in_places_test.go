package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
)

func TestNewAllTransitionPlacesInPlaces(t *testing.T) {
	assert.Equal(t, &AllTransitionPlacesInPlaces{}, NewAllTransitionPlacesInPlaces())
}

func TestAllTransitionPlacesInPlaces_Validate_ValidCfg_NoErr(t *testing.T) {
	minCfg := cfg.Minimal{
		Start:  "",
		Finish: "",
		Places: []cfg.StringID{"a"},
		Transitions: map[string]cfg.MinimalTransition{
			"b": {
				From: []cfg.StringID{"a"},
			},
		},
	}

	v := NewAllTransitionPlacesInPlaces()

	assert.NoError(t, v.Validate(minCfg))
}

func TestAllTransitionPlacesInPlaces_Validate_NotValidCfg_ExpectedErr(t *testing.T) {
	minCfg := cfg.Minimal{
		Start:  "",
		Finish: "",
		Places: nil,
		Transitions: map[string]cfg.MinimalTransition{
			"b": {
				From: []cfg.StringID{"a"},
			},
		},
	}

	v := NewAllTransitionPlacesInPlaces()

	assert.Equal(
		t,
		NewError().Addf("place with id 'a' from transitions is not found in places"),
		v.Validate(minCfg),
	)
}
