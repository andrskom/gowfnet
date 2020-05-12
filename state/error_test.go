package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testingErrCode1 ErrCode = "testing.errCode1"
const testingErrCode2 ErrCode = "testing.errCode2"

func TestNewError(t *testing.T) {
	require.Equal(
		t,
		&Error{
			code:    ErrCodeUnknown,
			message: "a",
		},
		NewError(ErrCodeUnknown, "a"),
	)
}

func TestNewErrorf(t *testing.T) {
	t.Run("no format", func(t *testing.T) {
		require.Equal(
			t,
			&Error{
				code:    ErrCodeUnknown,
				message: "a",
			},
			NewErrorf(ErrCodeUnknown, "a"),
		)
	})
	t.Run("with format", func(t *testing.T) {
		require.Equal(
			t,
			&Error{
				code:    ErrCodeUnknown,
				message: "a b",
			},
			NewErrorf(ErrCodeUnknown, "a %s", "b"),
		)
	})
}

func TestError_Is(t *testing.T) {
	t.Run("not is errStack", func(t *testing.T) {
		require.False(t, NewErrorf(ErrCodeUnknown, "a").Is(ErrCodeNetDoesntKnowAboutPlace))
	})
	t.Run("is errStack", func(t *testing.T) {
		require.True(t, NewErrorf(ErrCodeUnknown, "a").Is(ErrCodeUnknown))
	})
}

func TestError_Error(t *testing.T) {
	require.Equal(
		t,
		"a",
		NewErrorf(ErrCodeUnknown, "a").Error(),
	)
}

func TestErrorIs(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.False(t, ErrorIs(ErrCodeUnknown, nil))
	})
	t.Run("not Error type", func(t *testing.T) {
		require.False(t, ErrorIs(ErrCodeUnknown, errors.New("a")))
	})
	t.Run("not expected code", func(t *testing.T) {
		require.False(t, ErrorIs(ErrCodeUnknown, NewError(ErrCodeNetDoesntKnowAboutPlace, "a")))
	})
	t.Run("expected code", func(t *testing.T) {
		require.True(t, ErrorIs(ErrCodeUnknown, NewError(ErrCodeUnknown, "a")))
	})
}

func TestBuildError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.Nil(t, BuildError(nil))
	})
	t.Run("Error type", func(t *testing.T) {
		err := NewError(ErrCodeNetDoesntKnowAboutPlace, "a")
		require.Equal(t, err, BuildError(err))
	})
	t.Run("errStack not Error type", func(t *testing.T) {
		require.Equal(t, NewError(ErrCodeUnknown, "a"), BuildError(errors.New("a")))
	})
}

func TestError_Serialization(t *testing.T) {
	errModel := NewError(ErrCodeStateAlreadyHasTokenInPlace, "some message")
	bytes, err := json.Marshal(errModel)
	require.NoError(t, err)
	require.Equal(t, `{"code":"gowfnet.state.AlreadyHasTokenInPlace","message":"some message"}`, string(bytes))

	var errNewModel Error

	require.NoError(t, json.Unmarshal(bytes, &errNewModel))
	require.Equal(t, errModel, &errNewModel)
}

func TestError_UnserializeFromBadData_UNmarshalErr(t *testing.T) {
	var errNewModel Error

	err := json.Unmarshal([]byte("[]"), &errNewModel)
	require.IsType(t, &json.UnmarshalTypeError{}, err)
}

func TestNewErrStack_DefaultData_ReturnsExpectedData(t *testing.T) {
	assert.Equal(t, make([]Error, 0), NewErrStack().GetErrs())
}

func TestErrStack_Add_Nil_ExpectedPanic(t *testing.T) {
	defer func() {
		r := recover()
		assert.Equal(t, "must use not nil err", fmt.Sprintf("%+v", r))
	}()

	NewErrStack().Add(nil)
}

func TestErrStack_Add_Err_ExpectedStateOfStack(t *testing.T) {
	err := NewError(testingErrCode1, "")

	stack := NewErrStack()
	stack.Add(err)

	assert.Len(t, stack.GetErrs(), 1)
	assert.Equal(t, *err, stack.GetErrs()[0])
}

func TestErrStack_HasErrs_EmptyStack_ReturnsFalse(t *testing.T) {
	stack := NewErrStack()

	assert.False(t, stack.HasErrs())
}

func TestErrStack_HasErrs_SetStack_ReturnsTrue(t *testing.T) {
	stack := NewErrStack()
	stack.Add(NewError(testingErrCode1, ""))

	assert.True(t, stack.HasErrs())
}

func TestErrStack_Error(t *testing.T) {
	stack := NewErrStack()
	stack.Add(NewError(testingErrCode1, "a"))
	stack.Add(NewError(testingErrCode2, "b"))

	assert.Equal(t, `0) a;
1) b;
`, stack.Error())
}

func TestErrStack_Serialization(t *testing.T) {
	stack := NewErrStack()
	stack.Add(NewError(testingErrCode1, "a"))
	stack.Add(NewError(testingErrCode2, "b"))

	data, err := json.Marshal(stack)
	assert.NoError(t, err)
	assert.Equal(
		t,
		"{\"stack\":[{\"code\":\"testing.errCode1\",\"message\":\"a\"},"+
			"{\"code\":\"testing.errCode2\",\"message\":\"b\"}]}",
		string(data),
	)

	var newErrStack ErrStack

	assert.NoError(t, json.Unmarshal(data, &newErrStack))
	assert.Equal(t, stack, &newErrStack)
}

func TestErrStack_UnserializeBadData_ExpectedErr(t *testing.T) {
	var newErrStack ErrStack

	assert.IsType(t, &json.UnmarshalTypeError{}, json.Unmarshal([]byte("[]"), &newErrStack))
}
