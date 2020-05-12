package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrskom/gowfnet/cfg"
)

func TestNewNotNil(t *testing.T) {
	assert.Equal(t, &NotNil{}, NewNotNil())
}

func TestNotNil_Validate_ValidCfg_NoErr(t *testing.T) {
	v := NewNotNil()

	assert.NoError(t, v.Validate(&cfg.Minimal{}))
}

func TestNotNil_Validate_NotValidCfg_ExpectedErr(t *testing.T) {
	v := NewNotNil()

	assert.Equal(t, BuildErrorf("config of net can't be nil"), v.Validate(nil))
}
