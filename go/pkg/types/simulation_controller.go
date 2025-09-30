package types

import (
	"time"
)

// SimulationController provides an interface for simulation handling for individual programs
type SimulationController interface {
	// InjectSatellites injects the satellites to simulation
	InjectSatellites([]Node) error

	// InjectGroundStations injects the ground stations to simulation
	InjectGroundStations([]Node) error

	// StartAutorun starts autorun and returns running chan struct
	StartAutorun() <-chan struct{}

	// StopAutorun stops autorun gracefully
	StopAutorun()

	// StepBySeconds add given amount of seconds to current simulation time.
	// Then calculating the simulation step for the new simulation time.
	StepBySeconds(seconds float64)

	// StepByTime set the new simulation time.
	// Then calculating the simulation step for the new simulation time.
	StepByTime(newTime time.Time)

	// GetAllNodes returns all nodes in this simulation (satellites and ground stations combined)
	GetAllNodes() []Node

	// GetSatellites returns all satellites in this simulation
	GetSatellites() []Satellite

	// GetGroundStations returns all ground stations in this simulation
	GetGroundStations() []GroundStation

	// GetSimulationTime return the current simulation time
	GetSimulationTime() time.Time

	// GetStatePluginRepository returns the repository which handles the state plugins
	GetStatePluginRepository() *StatePluginRepository

	// Close gracefully shutdown the simulation
	Close()
}
