package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
)

func TestNewFinishPlaceInPlaces(t *testing.T) {
	assert.Equal(t, &FinishPlaceInPlaces{}, NewFinishPlaceInPlaces())
}

func TestFinishPlaceInPlaces_Validate_ValidCfg_NoErr(t *testing.T) {
	minCfg := &cfg.Minimal{
		Start:       "",
		Finish:      "a",
		Places:      []cfg.StringID{"a"},
		Transitions: nil,
	}

	v := NewFinishPlaceInPlaces()

	assert.NoError(t, v.Validate(minCfg))
}

func TestFinishPlaceInPlaces_Validate_NotValidCfg_ExpectedErr(t *testing.T) {
	minCfg := &cfg.Minimal{
		Start:       "",
		Finish:      "a",
		Places:      nil,
		Transitions: nil,
	}

	v := NewFinishPlaceInPlaces()

	assert.Equal(t, BuildErrorf("finish place is not found in places"), v.Validate(minCfg))
}
