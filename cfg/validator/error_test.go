package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	require.Equal(
		t,
		&Error{
			list: make([]string, 0),
		},
		NewError(),
	)
}

func TestError_Has(t *testing.T) {
	err := NewError()
	require.False(t, err.Has())
	err.Addf("blah")
	require.True(t, err.Has())
}

func TestError_Add(t *testing.T) {
	err := NewError()
	err.Addf("a")
	err.Addf("b %s", "c")
	require.Equal(
		t,
		&Error{
			list: []string{
				"a",
				"b c",
			},
		},
		err,
	)
}

func TestError_Error(t *testing.T) {
	err := NewError()
	err.Addf("a")
	err.Addf("b %s", "c")
	require.Equal(
		t,
		` - a
 - b c
`,
		err.Error(),
	)
}

func TestError_Get(t *testing.T) {
	err := NewError()
	err.Addf("a")
	err.Addf("b %s", "c")
	require.Equal(
		t,
		[]string{"a", "b c"},
		err.Get(),
	)
}

func TestPrepareResultErr_HasNotErr_Nil(t *testing.T) {
	err := NewError()
	assert.Nil(t, PrepareResultErr(err))
}

func TestPrepareResultErr_HasErr_Err(t *testing.T) {
	err := NewError()
	err.Addf("a")
	assert.Equal(t, &Error{list: []string{"a"}}, PrepareResultErr(err))
}

func TestBuildErrorf(t *testing.T) {
	err := BuildErrorf("%s", "a")
	assert.Equal(t, " - a\n", err.Error())
}
