package gowfnet

import (
	"fmt"
	"strings"
)

// Cfg is config for wf network.
type Cfg struct {
	Start       string                   `json:"start"`
	Finish      string                   `json:"finish"`
	Places      []string                 `json:"places"`
	Transitions map[string]CfgTransition `json:"transitions"`
}

// CfgTransition is struct of transition.
//
// If you want to use IsAutomatic, add AutomaticListenerMiddleware to listener,
// set is automatic and if transition is possible, it will be done.
type CfgTransition struct {
	From        []string `json:"from"`
	To          []string `json:"to"`
	IsAutomatic bool     `json:"isAutomatic,omitempty"`
}

// CfgValidateError is struct for validation errors.
type CfgValidateError struct {
	errs []string
}

// NewCfgValidateError init error
func NewCfgValidateError() *CfgValidateError {
	return &CfgValidateError{
		errs: make([]string, 0),
	}
}

// Addf message to err list.
func (c *CfgValidateError) Addf(format string, args ...interface{}) {
	c.errs = append(c.errs, fmt.Sprintf(format, args...))
}

// Has errors.
func (c *CfgValidateError) Has() bool {
	return len(c.errs) > 0
}

func (c *CfgValidateError) Error() string {
	return " - " + strings.Join(c.errs, "\n - ") + "\n"
}

// Validate config.
// nolint:funlen
func (c *Cfg) Validate() error {
	err := NewCfgValidateError()

	if len(c.Start) == 0 {
		err.Addf("Start place is empty")
	}

	if len(c.Finish) == 0 {
		err.Addf("Finish place is empty")
	}

	if len(c.Places) == 0 {
		err.Addf("Places is empty")
	}

	if len(c.Transitions) == 0 {
		err.Addf("Transitions is empty")
	}

	if err.Has() {
		return err
	}

	if c.Start == c.Finish {
		err.Addf("Start place can't be equal finish place")
	}

	if err.Has() {
		return err
	}

	placesRegistry := make(map[string]bool)
	hasStartPlace := false
	hasFinishPlace := false

	for _, place := range c.Places {
		if place == c.Start {
			hasStartPlace = true
		}

		if place == c.Finish {
			hasFinishPlace = true
		}

		if _, ok := placesRegistry[place]; ok {
			err.Addf("Place '%s' met two or more times in places", place)
		}

		placesRegistry[place] = false
	}

	if !hasStartPlace {
		err.Addf("Start place '%s' is not in places list '%+v'", c.Start, c.Places)
	}

	if !hasFinishPlace {
		err.Addf("Finish place '%s' is not in places list '%+v'", c.Finish, c.Places)
	}

	// nolint:godox
	// TODO: graph validation, try to find dead places, not connected by transition parts of net and dead cycles.

	if err.Has() {
		return err
	}

	return nil
}
