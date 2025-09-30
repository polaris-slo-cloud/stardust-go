package types

// StatePluginBuilder provides the interface for builder of StatePlugin
type StatePluginBuilder interface {
	BuildPlugins(pluginNames []string) ([]StatePlugin, error)
}
