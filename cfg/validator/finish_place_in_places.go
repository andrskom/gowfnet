package validator

import "github.com/andrskom/gowfnet/cfg"

type FinishPlaceInPlaces struct {
}

func NewFinishPlaceInPlaces() *FinishPlaceInPlaces {
	return &FinishPlaceInPlaces{}
}

func (s *FinishPlaceInPlaces) Validate(c cfg.Interface) error {
	for _, place := range c.GetPlaces() {
		if place.GetID() == c.GetFinish().GetID() {
			return nil
		}
	}

	err := NewError()
	err.Addf("finish place is not found in places")

	return err
}
