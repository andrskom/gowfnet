package validator

import (
	"github.com/andrskom/gowfnet/cfg"
)

type Validator interface {
	Validate(c cfg.Interface) error
}

type Combined struct {
	validators []Validator
}

func New(validators ...Validator) *Combined {
	if validators == nil {
		validators = make([]Validator, 0)
	}

	return &Combined{
		validators: validators,
	}
}

func (c *Combined) Validate(cfg cfg.Interface) error {
	for _, v := range c.validators {
		if err := v.Validate(cfg); err != nil {
			return err
		}
	}

	return nil
}

// NewCombinedWithAllValidators return component with all validators.
func NewCombinedWithAllValidators() *Combined {
	return New(
		NewNotNil(),
		NewEmpty(),
		NewStartPlaceInPlaces(),
		NewFinishPlaceInPlaces(),
		NewAllPlacesInTransitions(),
		NewAllTransitionPlacesInPlaces(),
		NewDuplicatedPlacesInPlaces(),
		NewDuplicatedPlacesInTransitions(),
		NewDeadPlaces(NewCfgTreeBuilder()),
		NewNonFinishPlaces(NewCfgTreeBuilder()),
	)
}
