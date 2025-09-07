package types

import "time"

type SimulationController interface {
	InjectSatellites([]Node) error
	InjectGroundStations([]Node) error
	StartAutorun() <-chan struct{}
	StopAutorun()
	StepBySeconds(seconds float64)
	StepByTime(newTime time.Time)
	GetAllNodes() []Node
	GetSatellites() []Node
	GetGroundStations() []Node
	GetSimulationTime() time.Time
}
