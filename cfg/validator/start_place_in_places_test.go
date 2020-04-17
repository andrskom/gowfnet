package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
)

func TestNewStartPlaceInPlaces(t *testing.T) {
	assert.Equal(t, &StartPlaceInPlaces{}, NewStartPlaceInPlaces())
}

func TestStartPlaceInPlaces_Validate_ValidCfg_NoErr(t *testing.T) {
	minCfg := &cfg.Minimal{
		Start:       "a",
		Finish:      "",
		Places:      []cfg.StringID{"a"},
		Transitions: nil,
	}

	v := NewStartPlaceInPlaces()

	assert.NoError(t, v.Validate(minCfg))
}

func TestStartPlaceInPlaces_Validate_NotValidCfg_ExpectedErr(t *testing.T) {
	minCfg := &cfg.Minimal{
		Start:       "a",
		Finish:      "",
		Places:      nil,
		Transitions: nil,
	}

	v := NewStartPlaceInPlaces()

	assert.Equal(t, NewError().Addf("start place is not found in places"), v.Validate(minCfg))
}
