package cfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andrskom/gowfnet"
)

func TestNewRegistry(t *testing.T) {
	assert.Equal(t, &Registry{data: make(map[string]ConfigInterface)}, NewRegistry())
}

func TestRegistry_AddWithName_NilConfig_ExpectedErrAndExpectedState(t *testing.T) {
	reg := NewRegistry()
	err := reg.AddWithName("a", nil)
	assert.True(t, gowfnet.ErrorIs(ErrCodeNilCfg, err), "expected err with code "+ErrCodeNilCfg)
	assert.Equal(t, &Registry{data: make(map[string]ConfigInterface)}, reg)
}

func TestRegistry_AddWithName_DoubleUseName_ExpectedErrAndExpectedState(t *testing.T) {
	reg := NewRegistry()
	err := reg.AddWithName("a", &Minimal{Start: "b"})
	require.NoError(t, err)

	err = reg.AddWithName("a", &Minimal{Start: "c"})
	assert.True(
		t,
		gowfnet.ErrorIs(ErrCodeCfgAlreadyRegistered, err),
		"expected err with code "+ErrCodeCfgAlreadyRegistered,
	)
	assert.Equal(t, &Registry{data: map[string]ConfigInterface{"a": &Minimal{Start: "b"}}}, reg)
}

func TestRegistry_AddWithName_CoupleCorrect_NoErrAndExpectedState(t *testing.T) {
	reg := NewRegistry()
	err := reg.AddWithName("a", &Minimal{Start: "b"})
	require.NoError(t, err)
	err = reg.AddWithName("c", &Minimal{Start: "d"})
	require.NoError(t, err)

	assert.Equal(
		t,
		&Registry{
			data: map[string]ConfigInterface{
				"a": &Minimal{Start: "b"},
				"c": &Minimal{Start: "d"},
			},
		},
		reg,
	)
}

func TestRegistry_GetByName_UnknownName_ExpectedErr(t *testing.T) {
	reg := NewRegistry()

	res, err := reg.GetByName("a")
	assert.Nil(t, res)
	assert.True(
		t,
		gowfnet.ErrorIs(ErrCodeCfgNotRegistered, err),
		"expected err with code "+ErrCodeCfgNotRegistered,
	)
}

func TestRegistry_GetByName_SetCorrectData_ExpectedRes(t *testing.T) {
	reg := NewRegistry()
	err := reg.AddWithName("a", &Minimal{Start: "b"})
	require.NoError(t, err)

	res, err := reg.GetByName("a")
	assert.NoError(t, err)
	assert.Equal(t, &Minimal{Start: "b"}, res)
}
