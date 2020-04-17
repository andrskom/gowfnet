package validator

import "github.com/andrskom/gowfnet/cfg"

func buildPlaceRegistryFromTransitions(tr cfg.TransitionRegistryInterface) map[string]struct{} {
	res := make(map[string]struct{})

	for _, transition := range tr.GetAsMap() {
		for _, place := range transition.GetFrom() {
			res[place.GetID()] = struct{}{}
		}

		for _, place := range transition.GetTo() {
			res[place.GetID()] = struct{}{}
		}
	}

	return res
}
