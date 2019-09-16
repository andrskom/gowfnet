package gowfnet

// AutomaticListenerMiddleware is a listener for automatic make transition.
type AutomaticListenerMiddleware struct {
	net      *Net
	listener Listener
}

// NewAutomaticListenerMiddleware init middleware.
func NewAutomaticListenerMiddleware(
	net *Net,
	listener Listener,
) *AutomaticListenerMiddleware {
	return &AutomaticListenerMiddleware{
		net:      net,
		listener: listener,
	}
}

// AfterPlaced is implementation of event listener interface.
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
