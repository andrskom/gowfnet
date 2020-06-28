package cfg

import (
	"github.com/andrskom/gowfnet/state"
)

const (
	ErrCodeNilCfg               state.ErrCode = "gowfnet.cfg.nilConfig"
	ErrCodeCfgAlreadyRegistered state.ErrCode = "gowfnet.cfg.alreadyRegistered"
	ErrCodeCfgNotRegistered     state.ErrCode = "gowfnet.cfg.notRegistered"
)

// Registry is a registry for config.
// Very often you provide more than one configured net for users.
type Registry struct {
	data map[string]Interface
}

// NewRegistry init empty registry.
func NewRegistry() *Registry {
	return &Registry{data: make(map[string]Interface)}
}

// AddWithName returns err if one of params will be unexpected.
func (r Registry) AddWithName(name string, cfg Interface) error {
	if cfg == nil {
		return state.NewError(ErrCodeNilCfg, "can't set nil config to registry")
	}

	if _, ok := r.data[name]; ok {
		return state.NewError(ErrCodeCfgAlreadyRegistered, "config with the same name is already registered")
	}

	r.data[name] = cfg

	return nil
}

// GetByName return Config or err if config was not registered.
func (r Registry) GetByName(name string) (Interface, error) {
	out, ok := r.data[name]
	if !ok {
		return nil, state.NewError(ErrCodeCfgNotRegistered, "config with this name was not registered")
	}

	return out, nil
}
