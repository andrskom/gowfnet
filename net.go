package gowfnet

// Listener is an interface of set place listener.
// If you want to ue many listeners, you should create own listener with list of listeners or use it as middlewares.
// You can see example in AutomaticListenerMiddleware.
type Listener interface {
	AfterPlaced(state *State, placeID string, subject interface{})
}

// AutomaticListenerMiddleware is a listener for automatic make transition.
type AutomaticListenerMiddleware struct {
	errorLogger func(msg string)
	net         *Net
	listener    Listener
}

// NewAutomaticListenerMiddleware init middleware.
func NewAutomaticListenerMiddleware(
	errorLogger func(msg string),
	net *Net,
	listener Listener,
) *AutomaticListenerMiddleware {
	return &AutomaticListenerMiddleware{
		errorLogger: errorLogger,
		net:         net,
		listener:    listener,
	}
}

func (l *AutomaticListenerMiddleware) AfterPlaced(state *State, placeID string, subject interface{}) {
	if l.listener != nil {
		l.listener.AfterPlaced(state, placeID, subject)
	}

	place, err := l.net.GetPlace(placeID)
	if err != nil {
		state.SetError(err)
		return
	}
	for _, transition := range place.GetToTransitions() {
		if transition.isAutomatic {
			if err := l.net.Transit(state, transition.id, subject); err != nil {
				if !ErrorIs(ErrCodeStateHasNotTokenInPlace, err) {
					state.SetError(err)
				}
			}
		}
	}
}

// Net is built net.
type Net struct {
	startPlace  *Place
	places      map[string]*Place
	transitions map[string]*Transition
	listener    Listener
}

// BuildFromConfig new net.
func BuildFromConfig(cfg Cfg) (*Net, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	placeRegistry := make(map[string]*Place)

	for _, placeID := range cfg.Places {
		placeRegistry[placeID] = newPlace(placeID, placeID == cfg.Finish)
	}

	wf := &Net{
		startPlace:  placeRegistry[cfg.Start],
		places:      placeRegistry,
		transitions: make(map[string]*Transition),
	}

	for transitionID, transitionCfg := range cfg.Transitions {
		transition := newTransition(transitionID, transitionCfg.IsAutomatic)

		for _, placeID := range transitionCfg.From {
			placeRegistry[placeID].addToTransitions(transition)
			transition.addFromPlace(placeRegistry[placeID])
		}

		for _, placeID := range transitionCfg.To {
			transition.addToPlace(placeRegistry[placeID])
		}
	}

	return wf, nil
}

// Start net for the state.
func (n *Net) Start(state *State, subject interface{}) error {
	if state.IsStarted() {
		return NewError(ErrCodeStateAlreadyStarted, "State already started in net")
	}

	return n.process(state, subject, []string{}, []string{n.startPlace.id})
}

// Transit state to new place by transit.
func (n *Net) Transit(state *State, transitionID string, subject interface{}) error {
	if !state.IsStarted() {
		return NewError(ErrCodeStateIsNotStarted, "Can't transit, state is not started")
	}

	transition, ok := n.transitions[transitionID]
	if !ok {
		return NewErrorf(
			ErrCodeNetDoesntKnowAboutTransition,
			"Net doesn't know about transition '%s'",
			transitionID,
		)
	}

	return n.process(state, subject, transition.fromPlaceIDs, transition.toPlaceIDs)
}

func (n *Net) process(state *State, subject interface{}, fromPlaces []string, toPlaces []string) error {
	if err := state.MoveTokensFromPlacesToPlaces(fromPlaces, toPlaces); err != nil {
		return err
	}
	for _, place := range toPlaces {
		n.listener.AfterPlaced(state, place, subject)
	}
	return nil
}

// GetPlace return copy of place.
func (n *Net) GetPlace(placeID string) (*Place, error) {
	place, ok := n.places[placeID]
	if !ok {
		return nil, NewErrorf(
			ErrCodeNetDoesntKnowAboutPlace,
			"Net doesn't know about place '%s'",
			placeID,
		)
	}
	copyPlace := *place

	return &copyPlace, nil
}
