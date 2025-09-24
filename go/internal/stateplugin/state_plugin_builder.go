package stateplugin

import (
	"fmt"

	"github.com/keniack/stardustGo/pkg/types"
)

type StatePluginBuilder struct {
}

// NewStatePluginBuilder creates a new instance of StatePluginBuilder
func NewStatePluginBuilder() *StatePluginBuilder {
	return &StatePluginBuilder{}
}

// BuildPlugins constructs plugin instances based on provided names
func (pb *StatePluginBuilder) BuildPlugins(pluginNames []string) ([]types.StatePlugin, error) {
	var plugins []types.StatePlugin
	for _, name := range pluginNames {
		switch name {
		case "DummyPlugin":
			plugins = append(plugins, NewDummySunStatePlugin())
		default:
			return nil, fmt.Errorf("unknown plugin: %s", name)
		}
	}
	return plugins, nil
}
