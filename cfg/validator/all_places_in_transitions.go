package validator

import (
	"github.com/andrskom/gowfnet/cfg"
)

type AllPlacesInTransitions struct {
}

func NewAllPlacesInTransitions() *AllPlacesInTransitions {
	return &AllPlacesInTransitions{}
}

func (a *AllPlacesInTransitions) Validate(c cfg.Interface) error {
	m := buildPlaceRegistryFromTransitions(c.GetTransitions())
	err := NewError()

	for _, place := range c.GetPlaces() {
		if _, ok := m[place.GetID()]; !ok {
			err.Addf("transitions don't use place with id '%s'", place.GetID())
		}
	}

	return PrepareResultErr(err)
}
