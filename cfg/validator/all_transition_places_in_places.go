package validator

import (
	"github.com/andrskom/gowfnet/cfg"
)

type AllTransitionPlacesInPlaces struct {
}

func NewAllTransitionPlacesInPlaces() *AllTransitionPlacesInPlaces {
	return &AllTransitionPlacesInPlaces{}
}

func (a *AllTransitionPlacesInPlaces) Validate(c cfg.Interface) error {
	m := buildPlaceRegistryFromTransitions(c.GetTransitions())

	for _, place := range c.GetPlaces() {
		delete(m, place.GetID())
	}

	if len(m) == 0 {
		return nil
	}

	err := NewError()
	for k := range m {
		err.Addf("place with id '%s' from transitions is not found in places", k)
	}

	return err
}
