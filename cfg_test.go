package gowfnet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCfgValidateError(t *testing.T) {
	require.Equal(
		t,
		&CfgValidateError{
			errs: make([]string, 0),
		},
		NewCfgValidateError(),
	)
}

func TestCfgValidateError_Has(t *testing.T) {
	err := NewCfgValidateError()
	require.False(t, err.Has())
	err.Add("blah")
	require.True(t, err.Has())
}

func TestCfgValidateError_Add(t *testing.T) {
	err := NewCfgValidateError()
	err.Add("a")
	err.Add("b %s", "c")
	require.Equal(
		t,
		&CfgValidateError{
			errs: []string{
				"a",
				"b c",
			},
		},
		err,
	)
}

func TestCfgValidateError_Error(t *testing.T) {
	err := NewCfgValidateError()
	err.Add("a")
	err.Add("b %s", "c")
	require.Equal(
		t,
		` - a
 - b c
`,
		err.Error(),
	)
}

func TestCfg_Validate_EmptyFields_Err(t *testing.T) {
	dp := map[string]Cfg{
		"nil places and transitions": {
			Start:       "",
			Finish:      "",
			Places:      nil,
			Transitions: nil,
		},
		"empty places and transitions": {
			Start:       "",
			Finish:      "",
			Places:      []string{},
			Transitions: map[string]CfgTransition{},
		},
	}

	for descr, data := range dp {
		t.Run(descr, func(t *testing.T) {
			require.Equal(
				t,
				&CfgValidateError{
					errs: []string{
						"Start place is empty",
						"Finish place is empty",
						"Places is empty",
						"Transitions is empty",
					},
				},
				data.Validate(),
			)
		})
	}
}

func TestCfg_Validate_StartequalFinish_Err(t *testing.T) {
	cfg := Cfg{
		Start:  "equal",
		Finish: "equal",
		Places: []string{"equal"},
		Transitions: map[string]CfgTransition{
			"equal": {},
		},
	}

	require.Equal(
		t,
		&CfgValidateError{
			errs: []string{
				"Start place can't be equal finish place",
			},
		},
		cfg.Validate(),
	)
}

func TestCfg_Validate_RepeatPlace_Err(t *testing.T) {
	cfg := Cfg{
		Start:  "blah",
		Finish: "blah1",
		Places: []string{
			"blah",
			"blah1",
			"repeat",
			"repeat",
		},
		Transitions: map[string]CfgTransition{
			"equal": {},
		},
	}

	require.Equal(
		t,
		&CfgValidateError{
			errs: []string{
				"Place 'repeat' met two or more times in places",
			},
		},
		cfg.Validate(),
	)
}

func TestCfg_Validate_StartAndFinishNotInPlaces_Err(t *testing.T) {
	cfg := Cfg{
		Start:  "blah",
		Finish: "blah1",
		Places: []string{
			"blah2",
			"blah3",
		},
		Transitions: map[string]CfgTransition{
			"equal": {},
		},
	}

	require.Equal(
		t,
		&CfgValidateError{
			errs: []string{
				"Start place 'blah' is not in places list '[blah2 blah3]'",
				"Finish place 'blah1' is not in places list '[blah2 blah3]'",
			},
		},
		cfg.Validate(),
	)
}


func TestCfg_Validate_Correct_NoErr(t *testing.T) {
	cfg := Cfg{
		Start:  "blah",
		Finish: "blah1",
		Places: []string{
			"blah",
			"blah1",
		},
		Transitions: map[string]CfgTransition{
			"transition": {
				From: []string{
					"blah",
				},
				To: []string{
					"blah1",
				},
			},
		},
	}

	require.NoError(t,cfg.Validate())
}
