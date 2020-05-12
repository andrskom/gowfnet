package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
)

func TestNewDuplicatedPlacesInPlaces(t *testing.T) {
	assert.Equal(t, &DuplicatedPlacesInPlaces{}, NewDuplicatedPlacesInPlaces())
}

func TestDuplicatedPlacesInPlaces_Validate_RepeatVal_ExpectedErr(t *testing.T) {
	minCfg := cfg.Minimal{
		Places: []cfg.StringID{
			"a",
			"a",
			"b",
		},
	}

	v := NewDuplicatedPlacesInPlaces()

	assert.Equal(t, BuildErrorf("place with id 'a' is duplicated"), v.Validate(minCfg))
}

func TestDuplicatedPlacesInPlaces_Validate_CorrectCfg_NoErr(t *testing.T) {
	minCfg := cfg.Minimal{
		Places: []cfg.StringID{
			"a",
			"b",
		},
	}

	v := NewDuplicatedPlacesInPlaces()

	assert.NoError(t, v.Validate(minCfg))
}
