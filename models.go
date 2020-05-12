package gowfnet

// Transition is a struct contains info about transition.
type Transition struct {
	id           string
	toPlaces     []*Place
	toPlaceIDs   []string
	fromPlaceIDs []string
	isAutomatic  bool
}

func newTransition(id string, isAutomatic bool) *Transition {
	return &Transition{
		id:           id,
		toPlaces:     make([]*Place, 0),
		toPlaceIDs:   make([]string, 0),
		fromPlaceIDs: make([]string, 0),
		isAutomatic:  isAutomatic,
	}
}

func (t *Transition) addToPlace(newPlace *Place) {
	t.toPlaces = append(t.toPlaces, newPlace)
	t.toPlaceIDs = append(t.toPlaceIDs, newPlace.id)
}

func (t *Transition) addFromPlace(newPlace *Place) {
	t.fromPlaceIDs = append(t.fromPlaceIDs, newPlace.id)
}

// Place is a struct for a place.
type Place struct {
	id            string
	toTransitions []*Transition
	isFinished    bool
}

// GetToTransitions return copy of toTransitions.
func (p *Place) GetToTransitions() []Transition {
	res := make([]Transition, 0, len(p.toTransitions))
	for _, transition := range p.toTransitions {
		res = append(res, *transition)
	}

	return res
}

func newPlace(id string, isFinished bool) *Place {
	return &Place{
		id:            id,
		toTransitions: make([]*Transition, 0),
		isFinished:    isFinished,
	}
}

func (p *Place) addToTransitions(newTransition *Transition) {
	p.toTransitions = append(p.toTransitions, newTransition)
}
