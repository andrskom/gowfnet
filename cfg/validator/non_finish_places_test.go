// nolint:funlen
package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
)

func TestNewNonFinishPlaces(t *testing.T) {
	assert.Equal(t, &NonFinishPlaces{treeBuilder: NewCfgTreeBuilder()}, NewNonFinishPlaces(NewCfgTreeBuilder()))
}

func TestNonFinishPlaces_Validate_BuilderErr_TheSameErr(t *testing.T) {
	builderMock := NewBuilderMock()

	v := NewNonFinishPlaces(builderMock)

	var (
		minCfg  *cfg.Minimal
		mockRes *Tree
	)

	eErr := errors.New("expectedErr")
	builderMock.On("Build", minCfg).Return(mockRes, eErr)

	assert.Same(t, eErr, v.Validate(minCfg))
}

func TestNonFinishPlaces_Validate_BadTree_ExpectedErr(t *testing.T) {
	builderMock := NewBuilderMock()

	v := NewNonFinishPlaces(builderMock)

	var minCfg *cfg.Minimal

	res := &Tree{
		startNodeID: "a",
		registry:    nil,
	}

	builderMock.On("Build", minCfg).Return(res, nil)

	assert.Equal(t, ErrNodeIsNotFound, v.Validate(minCfg))
}

func TestNonFinishPlaces_Validate_NotValidCfg_ExpectedErr(t *testing.T) {
	type data struct {
		cfg      cfg.Interface
		expected *Error
	}

	dp := map[string]data{
		"simple net": {
			cfg: cfg.Minimal{
				Start:  "a",
				Finish: "c",
				Places: []cfg.StringID{"a", "b", "c"},
				Transitions: map[string]cfg.MinimalTransition{
					"d": {
						From: []cfg.StringID{"b"},
						To:   []cfg.StringID{"c"},
					},
				},
			},
			expected: BuildErrorf("place with id 'a' is non-finish place"),
		},
		"net with branch error": {
			cfg: cfg.Minimal{
				Start:  "a",
				Finish: "z",
				Places: []cfg.StringID{"a", "b", "z", "c"},
				Transitions: map[string]cfg.MinimalTransition{
					"t1": {
						From: []cfg.StringID{"a"},
						To:   []cfg.StringID{"b", "c"},
					},
					"t2": {
						From: []cfg.StringID{"b"},
						To:   []cfg.StringID{"z"},
					},
				},
			},
			expected: BuildErrorf("place with id 'c' is non-finish place"),
		},
	}

	v := NewNonFinishPlaces(NewCfgTreeBuilder())

	for desc, d := range dp {
		t.Run(desc, func(t *testing.T) {
			assert.Equal(t, d.expected, v.Validate(d.cfg)) // nolint:scopelint
		})
	}
}

func TestNonFinishPlaces_Validate_ValidCfg_NoErr(t *testing.T) {
	type data struct {
		cfg cfg.Interface
	}

	dp := map[string]data{
		"simple net": {
			cfg: cfg.Minimal{
				Start:  "a",
				Finish: "b",
				Places: []cfg.StringID{"a", "b"},
				Transitions: map[string]cfg.MinimalTransition{
					"d": {
						From: []cfg.StringID{"a"},
						To:   []cfg.StringID{"b"},
					},
				},
			},
		},
		"net with place not followed to finish": {
			cfg: cfg.Minimal{
				Start:  "a",
				Finish: "z",
				Places: []cfg.StringID{"a", "b", "z", "c"},
				Transitions: map[string]cfg.MinimalTransition{
					"t1": {
						From: []cfg.StringID{"a"},
						To:   []cfg.StringID{"b"},
					},
					"t2": {
						From: []cfg.StringID{"b", "c"},
						To:   []cfg.StringID{"z"},
					},
				},
			},
		},
		"net with cycle": {
			cfg: cfg.Minimal{
				Start:  "a",
				Finish: "z",
				Places: []cfg.StringID{"a", "b", "z", "c"},
				Transitions: map[string]cfg.MinimalTransition{
					"t1": {
						From: []cfg.StringID{"a"},
						To:   []cfg.StringID{"b"},
					},
					"t2": {
						From: []cfg.StringID{"b"},
						To:   []cfg.StringID{"c"},
					},
					"t3": {
						From: []cfg.StringID{"c"},
						To:   []cfg.StringID{"b"},
					},
					"t4": {
						From: []cfg.StringID{"b"},
						To:   []cfg.StringID{"z"},
					},
				},
			},
		},
	}

	v := NewNonFinishPlaces(NewCfgTreeBuilder())

	for desc, d := range dp {
		t.Run(desc, func(t *testing.T) {
			assert.NoError(t, v.Validate(d.cfg)) // nolint:scopelint
		})
	}
}
