package types

type ISimulationController interface {
	InjectSatellites([]Node) error
}
