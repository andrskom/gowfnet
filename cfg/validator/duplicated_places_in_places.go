package validator

import (
	"github.com/andrskom/gowfnet/cfg"
)

type DuplicatedPlacesInPlaces struct {
}

func NewDuplicatedPlacesInPlaces() *DuplicatedPlacesInPlaces {
	return &DuplicatedPlacesInPlaces{}
}

func (s *DuplicatedPlacesInPlaces) Validate(c cfg.Interface) error {
	err := NewError()
	checkMap := make(map[string]struct{})

	for _, place := range c.GetPlaces() {
		if _, ok := checkMap[place.GetID()]; ok {
			err.Addf("place with id '%s' is duplicated", place.GetID())
		}

		checkMap[place.GetID()] = struct{}{}
	}

	return PrepareResultErr(err)
}
