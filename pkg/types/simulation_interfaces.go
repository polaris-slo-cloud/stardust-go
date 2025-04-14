package types

type ISimulationController interface {
	InjectSatellites([]*INode) error
}
