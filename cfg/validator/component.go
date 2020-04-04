package validator

import (
	"github.com/andrskom/gowfnet/cfg"
)

type Validator interface {
	Validate(c cfg.Interface) error
}

type Component struct {
	validators []Validator
}

func New(validators ...Validator) *Component {
	if validators == nil {
		validators = make([]Validator, 0)
	}

	return &Component{
		validators: validators,
	}
}

func (c *Component) Validate(cfg cfg.Interface) error {
	for _, v := range c.validators {
		if err := v.Validate(cfg); err != nil {
			return err
		}
	}

	return nil
}
