package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
)

func TestNewAllPlacesInTransitions(t *testing.T) {
	assert.Equal(t, &AllPlacesInTransitions{}, NewAllPlacesInTransitions())
}

func TestAllPlacesInTransitions_Validate_ValidCfg_NoErr(t *testing.T) {
	minCfg := &cfg.Minimal{
		Start:  "",
		Finish: "",
		Places: []cfg.StringID{"a"},
		Transitions: map[string]cfg.MinimalTransition{
			"b": {
				To: []cfg.StringID{"a"},
			},
		},
	}

	v := NewAllPlacesInTransitions()

	assert.NoError(t, v.Validate(minCfg))
}

func TestAllPlacesInTransitions_Validate_NotValidCfg_ExpectedErr(t *testing.T) {
	minCfg := &cfg.Minimal{
		Start:       "",
		Finish:      "",
		Places:      []cfg.StringID{"a"},
		Transitions: nil,
	}

	v := NewAllPlacesInTransitions()

	assert.Equal(t, NewError().Addf("transitions don't use place with id 'a'"), v.Validate(minCfg))
}
