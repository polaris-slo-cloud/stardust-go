package simplugin

import (
	"fmt"

	"github.com/keniack/stardustGo/pkg/types"
)

type SimPluginBuilder struct {
}

// NewPluginBuilder creates a new instance of PluginBuilder
func NewPluginBuilder() *SimPluginBuilder {
	return &SimPluginBuilder{}
}

// BuildPlugins constructs plugin instances based on provided names
func (pb *SimPluginBuilder) BuildPlugins(pluginNames []string) ([]types.SimulationPlugin, error) {
	var plugins []types.SimulationPlugin
	for _, name := range pluginNames {
		switch name {
		case "DummyPlugin":
			plugins = append(plugins, &DummyPlugin{})
		default:
			return nil, fmt.Errorf("unknown plugin: %s", name)
		}
	}
	return plugins, nil
}
