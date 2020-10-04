package validator

import (
	"github.com/andrskom/gowfnet/cfg"
)

type DuplicatedPlacesInTransitions struct {
}

func NewDuplicatedPlacesInTransitions() *DuplicatedPlacesInTransitions {
	return &DuplicatedPlacesInTransitions{}
}

func (d *DuplicatedPlacesInTransitions) Validate(c cfg.Interface) error {
	err := NewError()

	for transitionName, transition := range c.GetTransitions().GetAsMap() {
		checkMap := make(map[string]struct{})

		for _, place := range transition.GetFrom() {
			if _, ok := checkMap[place.GetID()]; ok {
				err.Addf(
					"place with id '%s' is duplicated in transition with id '%s' in section from",
					place.GetID(), transitionName,
				)
			}

			checkMap[place.GetID()] = struct{}{}
		}

		checkMap = make(map[string]struct{})

		for _, place := range transition.GetTo() {
			if _, ok := checkMap[place.GetID()]; ok {
				err.Addf(
					"place with id '%s' is duplicated in transition with id '%s' in section to",
					place.GetID(), transitionName,
				)
			}

			checkMap[place.GetID()] = struct{}{}
		}
	}

	return PrepareResultErr(err)
}
