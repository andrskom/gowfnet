package cfg

import (
	"encoding/json"

	"github.com/andrskom/gowfnet"
)

const (
	ErrCodeUnknownTransitionID gowfnet.ErrCode = "gowfnet.cfg.unknownTransition"
)

// StringID is simple implementation of IDGetter interface.
type StringID string

// CreateStringID build StringID from built-int string.
func CreateStringID(data string) StringID {
	return StringID(data)
}

// GetID return value as a string.
func (s StringID) GetID() string {
	return string(s)
}

// Minimal is a implementation of Interface.
// This contains only required fields.
// The easiest ways of setting config are via const or via json.
type Minimal struct {
	Start       StringID                  `json:"start"`
	Finish      StringID                  `json:"finish"`
	Places      []StringID                `json:"places"`
	Transitions MinimalTransitionRegistry `json:"transitions"`
}

func (m Minimal) GetStart() IDGetter {
	return m.Start
}

func (m Minimal) GetFinish() IDGetter {
	return m.Finish
}

func (m Minimal) GetPlaces() []IDGetter {
	return convertSliceFromStringToInterface(m.Places)
}

func (m Minimal) GetTransitions() TransitionRegistryInterface {
	return m.Transitions
}

// MinimalTransition is a simple implementation of TransitionInterface.
// This contains only required fields.
type MinimalTransition struct {
	To   []StringID `json:"to"`
	From []StringID `json:"from"`
}

func (m MinimalTransition) GetFrom() []IDGetter {
	return convertSliceFromStringToInterface(m.From)
}

func (m MinimalTransition) GetTo() []IDGetter {
	return convertSliceFromStringToInterface(m.To)
}

// MinimalTransitionRegistry is a simple implementation of TransitionRegistryInterface.
// This contains only required fields.
type MinimalTransitionRegistry struct {
	data map[string]MinimalTransition
}

func (m *MinimalTransitionRegistry) UnmarshalJSON(bytes []byte) error {
	return json.Unmarshal(bytes, &m.data)
}

func (m MinimalTransitionRegistry) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

func (m MinimalTransitionRegistry) GetAsMap() map[string]TransitionInterface {
	out := make(map[string]TransitionInterface)
	for k, v := range m.data {
		out[k] = v
	}

	return out
}

func (m MinimalTransitionRegistry) GetByID(transitionID IDGetter) (TransitionInterface, error) {
	if _, ok := m.data[transitionID.GetID()]; !ok {
		return nil, gowfnet.NewError(ErrCodeUnknownTransitionID, "can't find transition for id in registry")
	}

	return m.data[transitionID.GetID()], nil
}

func convertSliceFromStringToInterface(in []StringID) []IDGetter {
	out := make([]IDGetter, 0, len(in))
	for _, v := range in {
		out = append(out, v)
	}

	return out
}
