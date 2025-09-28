package stateplugin

import (
	"fmt"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.StatePluginBuilder = (*DefaultStatePluginPrecompBuilder)(nil)

type DefaultStatePluginPrecompBuilder struct {
	filename string
}

// NewStatePluginPrecompBuilder creates a new instance of StatePluginPrecompBuilder
func NewStatePluginPrecompBuilder(filename string) *DefaultStatePluginPrecompBuilder {
	return &DefaultStatePluginPrecompBuilder{
		filename: filename,
	}
}

// BuildPlugins constructs plugin instances based on provided names
func (pb *DefaultStatePluginPrecompBuilder) BuildPlugins(pluginNames []string) ([]types.StatePlugin, error) {
	var plugins []types.StatePlugin
	for _, name := range pluginNames {
		switch name {
		case "DummyPlugin":
			plugins = append(plugins, NewDummySunStatePrecompPlugin(pb.filename))
		default:
			return nil, fmt.Errorf("unknown plugin: %s", name)
		}
	}
	return plugins, nil
}
