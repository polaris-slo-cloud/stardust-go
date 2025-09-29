package types

type StatePluginBuilder interface {
	BuildPlugins(pluginNames []string) ([]StatePlugin, error)
}
