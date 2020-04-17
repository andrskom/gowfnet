package validator

import (
	"github.com/andrskom/gowfnet/cfg"
)

type Empty struct {
}

func NewEmpty() *Empty {
	return &Empty{}
}

func (e Empty) Validate(c cfg.Interface) error {
	err := NewError()

	if len(c.GetStart().GetID()) == 0 {
		err.Addf("start place id is empty")
	}

	if len(c.GetFinish().GetID()) == 0 {
		err.Addf("finish place id is empty")
	}

	if len(c.GetPlaces()) == 0 {
		err.Addf("places is empty")
	}

	if len(c.GetTransitions().GetAsMap()) == 0 {
		err.Addf("transitions registry is empty")
	}

	return PrepareResultErr(err)
}
