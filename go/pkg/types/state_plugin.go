package types

import "reflect"

type StatePlugin interface {
	PostSimulationStep(simulationController SimulationController)
}

type StatePluginRepository struct {
	plugins map[reflect.Type]StatePlugin
}

// NewStatePluginRepository creates a new StatePluginRepository and initializes it with the provided plugins.
func NewStatePluginRepository(plugins []StatePlugin) *StatePluginRepository {
	repo := &StatePluginRepository{
		plugins: make(map[reflect.Type]StatePlugin),
	}
	for _, plugin := range plugins {
		// Use the concrete type of the plugin as the key
		typ := reflect.TypeOf(plugin)
		repo.plugins[typ] = plugin
	}
	return repo
}

// GetAllPlugins returns all registered plugins.
func (r *StatePluginRepository) GetAllPlugins() []StatePlugin {
	plugins := make([]StatePlugin, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// GetStatePlugin is a generic function that retrieves a plugin of type T from the repository.
// It panics if the plugin is not found or if the type assertion fails.
func GetStatePlugin[T StatePlugin](r *StatePluginRepository) T {
	typ := reflect.TypeOf(*new(T))
	plugin, ok := r.plugins[typ]
	if !ok {
		panic("plugin not found")
	}
	return plugin.(T)
}
