package cfg

type IDGetter interface {
	GetID() string
}

type TransitionInterface interface {
	GetFrom() []IDGetter
	GetTo() []IDGetter
}

type TransitionRegistryInterface interface {
	GetAsMap() map[string]TransitionInterface
	GetByID(transitionID IDGetter) (TransitionInterface, error)
}

type Interface interface {
	GetStart() IDGetter
	GetFinish() IDGetter
	GetPlaces() []IDGetter
	GetTransitions() TransitionRegistryInterface
}
