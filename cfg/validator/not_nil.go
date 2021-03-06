package validator

import (
	"github.com/andrskom/gowfnet/cfg"
)

type NotNil struct {
}

func NewNotNil() *NotNil {
	return &NotNil{}
}

func (n *NotNil) Validate(c cfg.Interface) error {
	if c != nil {
		return nil
	}

	err := NewError()
	err.Addf("config of net can't be nil")

	return err
}
