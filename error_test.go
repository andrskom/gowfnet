package gowfnet

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

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
	t.Run("not is error", func(t *testing.T) {
		require.False(t, NewErrorf(ErrCodeUnknown, "a").Is(ErrCodeNetDoesntKnowAboutPlace))
	})
	t.Run("is error", func(t *testing.T) {
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
	t.Run("error not Error type", func(t *testing.T) {
		require.Equal(t, NewError(ErrCodeUnknown, "a"), BuildError(errors.New("a")))
	})
}
