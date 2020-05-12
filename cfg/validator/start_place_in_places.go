package validator

import "github.com/andrskom/gowfnet/cfg"

type StartPlaceInPlaces struct {
}

func NewStartPlaceInPlaces() *StartPlaceInPlaces {
	return &StartPlaceInPlaces{}
}

func (s *StartPlaceInPlaces) Validate(c cfg.Interface) error {
	for _, place := range c.GetPlaces() {
		if place.GetID() == c.GetStart().GetID() {
			return nil
		}
	}

	err := NewError()
	err.Addf("start place is not found in places")

	return err
}
