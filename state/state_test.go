package state

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewState(t *testing.T) {
	state := NewState()
	assert.Equal(t, make([]string, 0), state.GetPlaces())
	assert.Equal(t, NewErrStack(), state.GetErrorStack())
	assert.False(t, state.IsFinished())
	assert.False(t, state.IsStarted())
	assert.False(t, state.IsError())
}

func TestState_GetPlaces(t *testing.T) {
	state := NewState()
	require.NoError(t, state.MoveTokensFromPlacesToPlaces([]string{}, []string{"a"}))
	assert.Equal(t, []string{"a"}, state.GetPlaces())
}

func TestState_AddError(t *testing.T) {
	state := NewState()

	{
		state.AddError(errors.New("a"))
		errStack := state.GetErrorStack()

		assert.Len(t, errStack.GetErrs(), 1)
		assert.Equal(t,
			Error{
				code:    ErrCodeUnknown,
				message: "a",
			},
			errStack.GetErrs()[0],
		)
	}

	{
		state.AddError(NewError(testingErrCode1, "b"))
		errStack := state.GetErrorStack()

		assert.Len(t, errStack.GetErrs(), 2)
		assert.Equal(t,
			Error{
				code:    ErrCodeUnknown,
				message: "a",
			},
			errStack.GetErrs()[0],
		)
		assert.Equal(t,
			Error{
				code:    testingErrCode1,
				message: "b",
			},
			errStack.GetErrs()[1],
		)
	}
}

func TestState_SetFinished(t *testing.T) {
	state := NewState()

	{
		err := state.SetFinished()
		assert.NoError(t, err)
		assert.True(t, state.IsFinished(), "State must be finished")
	}

	{
		err := state.SetFinished()
		assert.Equal(
			t,
			&Error{
				code:    ErrCodeStateIsAlreadyFinished,
				message: "Can't set finished state, because state is already finished",
			},
			err,
		)
		assert.True(t, state.IsFinished(), "State must be finished")
	}
}

func TestState_MoveTokensFromPlacesToPlaces_ErrState_ExpectedErr(t *testing.T) {
	state := NewState()
	state.AddError(errors.New("a"))

	err := state.MoveTokensFromPlacesToPlaces([]string{}, []string{})
	assert.Equal(
		t,
		&Error{code: ErrCodeStateIsErrorState, message: "Can't process state to new places, state is errStack"},
		err,
	)
}

func TestState_MoveTokensFromPlacesToPlaces_StateIsFinished_ExpectedErr(t *testing.T) {
	state := NewState()
	require.NoError(t, state.SetFinished())

	err := state.MoveTokensFromPlacesToPlaces([]string{}, []string{})
	assert.Equal(
		t,
		&Error{code: ErrCodeStateIsFinished, message: "Can't process state to new places, state is finished"},
		err,
	)
}

func TestState_MoveTokensFromPlacesToPlaces_FromPLaceWithoutToken_ExpectedErr(t *testing.T) {
	state := NewState()

	err := state.MoveTokensFromPlacesToPlaces([]string{"a"}, []string{})
	assert.Equal(
		t,
		&Error{code: ErrCodeStateHasNotTokenInPlace, message: "State has not token in place 'a', state places: map[]"},
		err,
	)
}

func TestState_MoveTokensFromPlacesToPlaces_ToPLaceAlreadyWithToken_ExpectedErr(t *testing.T) {
	state := NewState()
	{
		err := state.MoveTokensFromPlacesToPlaces([]string{}, []string{"a", "b"})
		require.NoError(t, err)
	}

	err := state.MoveTokensFromPlacesToPlaces([]string{"a"}, []string{"b"})
	assert.Equal(
		t,
		&Error{
			code:    ErrCodeStateAlreadyHasTokenInPlace,
			message: "State already has token in place 'b', state places: map[a:{} b:{}]",
		},
		err,
	)
}

func TestState_MoveTokensFromPlacesToPlaces_CorrectStateForOperation_ExpectedState(t *testing.T) {
	state := NewState()
	{
		err := state.MoveTokensFromPlacesToPlaces([]string{}, []string{"a", "b"})
		require.NoError(t, err)
	}

	err := state.MoveTokensFromPlacesToPlaces([]string{"a", "b"}, []string{"c", "d"})
	assert.NoError(t, err)
	assert.Len(t, state.GetPlaces(), 2)
	assert.Contains(t, state.GetPlaces(), "c")
	assert.Contains(t, state.GetPlaces(), "d")
}

func TestState_Serialization(t *testing.T) {
	state := NewState()
	require.NoError(t, state.MoveTokensFromPlacesToPlaces([]string{}, []string{"a"}))
	state.AddError(errors.New("b"))
	require.NoError(t, state.SetFinished())

	bytes, err := json.Marshal(state)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"{\"places\":[\"a\"],\"errStack\":{\"stack\":[{\"code\":\"gowfnet.unknown\",\"message\":\"b\"}]},"+
			"\"isFinished\":true}",
		string(bytes),
	)

	{
		var newState State
		err := json.Unmarshal(bytes, &newState)
		assert.NoError(t, err)
		assert.Equal(t, state, &newState)
	}
}

func TestState_UnmarshalJSON_UnexpectedJSON_ExpectedErr(t *testing.T) {
	var state State
	err := json.Unmarshal([]byte("[]"), &state)
	assert.IsType(t, &json.UnmarshalTypeError{}, err)
}
