package gowfnet

// Registry is a net registry.
type Registry struct {
	data map[string]*Net
}

// NewRegistry init registry.
func NewRegistry() *Registry {
	return &Registry{data: make(map[string]*Net)}
}

// BuildRegistryFromCfgMap init registry from map of cfg.
func BuildRegistryFromCfgMap(cfgMap map[string]Cfg) (*Registry, error) {
	registry := NewRegistry()

	for name, cfg := range cfgMap {
		net, err := BuildFromConfig(cfg)
		if err != nil {
			return nil, err
		}
		if err := registry.Add(name, net); err != nil {
			return nil, err
		}
	}

	return registry, nil
}

// Add new net.
func (r *Registry) Add(name string, net *Net) error {
	if _, ok := r.data[name]; ok {
		return NewErrorf(
			ErrCodeRegistryNetAlreadyRegistered,
			"Net with name '%s' already registered",
			name,
		)
	}
	r.data[name] = net

	return nil
}

// Get return net if exists.
func (r *Registry) Get(name string) (*Net, error) {
	res, ok := r.data[name]
	if !ok {
		return nil, NewErrorf(
			ErrCodeRegistryNetNotRegistered,
			"Net with name '%s' is not registered",
			name,
		)
	}

	return res, nil
}
