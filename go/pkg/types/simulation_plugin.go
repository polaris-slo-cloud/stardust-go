package types

// SimulationPlugin defines the interface for simulation plugins.
type SimulationPlugin interface {
	// Name returns the name of the simulation plugin
	Name() string

	// PostSimulationStep is called by SimulationController after a step is completed.
	// Please provide your SimulationPlugin Logic in this method.
	PostSimulationStep(simualtion SimulationController) error
}
