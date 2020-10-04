package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
)

func TestNewEmpty(t *testing.T) {
	assert.Equal(t, &Empty{}, NewEmpty())
}

func TestEmpty_Validate_ValidCfg_NoErr(t *testing.T) {
	v := NewEmpty()

	minCfg := cfg.Minimal{
		Start:       "a",
		Finish:      "b",
		Places:      []cfg.StringID{"c"},
		Transitions: cfg.MinimalTransitionRegistry{"d": {}},
	}

	assert.NoError(t, v.Validate(minCfg))
}

func TestEmpty_Validate_NotValidCfg_ExpectedErr(t *testing.T) {
	type Data struct {
		cfg         cfg.Minimal
		expectedErr error
	}

	dp := map[string]Data{
		"empty start": {
			cfg: cfg.Minimal{
				Start:       "",
				Finish:      "b",
				Places:      []cfg.StringID{"c"},
				Transitions: cfg.MinimalTransitionRegistry{"d": {}},
			},
			expectedErr: BuildErrorf("start place id is empty"),
		},
		"empty finish": {
			cfg: cfg.Minimal{
				Start:       "a",
				Finish:      "",
				Places:      []cfg.StringID{"c"},
				Transitions: cfg.MinimalTransitionRegistry{"d": {}},
			},
			expectedErr: BuildErrorf("finish place id is empty"),
		},
		"empty places": {
			cfg: cfg.Minimal{
				Start:       "a",
				Finish:      "b",
				Places:      nil,
				Transitions: cfg.MinimalTransitionRegistry{"d": {}},
			},
			expectedErr: BuildErrorf("places is empty"),
		},
		"empty transitions": {
			cfg: cfg.Minimal{
				Start:       "a",
				Finish:      "b",
				Places:      []cfg.StringID{"c"},
				Transitions: nil,
			},
			expectedErr: BuildErrorf("transitions registry is empty"),
		},
	}

	v := NewEmpty()

	for desc, data := range dp {
		t.Run(desc, func(t *testing.T) {
			assert.Equal(t, data.expectedErr, v.Validate(data.cfg)) // nolint:scopelint
		})
	}
}
