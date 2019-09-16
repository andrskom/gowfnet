package gowfnet

// AutomaticListenerMiddleware is a listener for automatic make transition.
type AutomaticListenerMiddleware struct {
	listener Listener
}

// NewAutomaticListenerMiddleware init middleware.
func NewAutomaticListenerMiddleware(
	listener Listener,
) *AutomaticListenerMiddleware {
	return &AutomaticListenerMiddleware{
		listener: listener,
	}
}

// AfterPlaced is implementation of event listener interface.
func (l *AutomaticListenerMiddleware) AfterPlaced(net *Net, state *State, placeID string, subject interface{}) {
	if l.listener != nil {
		l.listener.AfterPlaced(net, state, placeID, subject)
	}

	place, err := net.GetPlace(placeID)
	if err != nil {
		state.SetError(err)
		return
	}
	for _, transition := range place.GetToTransitions() {
		if transition.isAutomatic {
			if err := net.Transit(state, transition.id, subject); err != nil {
				if !ErrorIs(ErrCodeStateHasNotTokenInPlace, err) {
					state.SetError(err)
				}
			}
		}
	}
}
