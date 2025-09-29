package stateplugin

import (
	"fmt"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.StatePluginBuilder = (*DefaultStatePluginBuilder)(nil)

type DefaultStatePluginBuilder struct {
}

// NewStatePluginBuilder creates a new instance of StatePluginBuilder
func NewStatePluginBuilder() *DefaultStatePluginBuilder {
	return &DefaultStatePluginBuilder{}
}

// BuildPlugins constructs plugin instances based on provided names
func (pb *DefaultStatePluginBuilder) BuildPlugins(pluginNames []string) ([]types.StatePlugin, error) {
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
