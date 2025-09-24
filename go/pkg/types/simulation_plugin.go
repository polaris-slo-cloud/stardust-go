package types

// SimulationPlugin defines the interface for simulation plugins.
type SimulationPlugin interface {
	Name() string
	PostSimulationStep(simualtion SimulationController) error
}
