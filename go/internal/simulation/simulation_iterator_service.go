package simulation

import (
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.SimulationController = (*SimulationIteratorService)(nil)

type SimulationIteratorService struct {
}

func (s *SimulationIteratorService) InjectSatellites([]types.Node) error {
	return nil
}
func (s *SimulationIteratorService) InjectGroundStations([]types.Node) error {
	return nil
}
func (s *SimulationIteratorService) StartAutorun() <-chan struct{} {
	return nil
}
func (s *SimulationIteratorService) StopAutorun() {

}
func (s *SimulationIteratorService) StepBySeconds(seconds float64) {

}
func (s *SimulationIteratorService) StepByTime(newTime time.Time) {

}
func (s *SimulationIteratorService) GetAllNodes() []types.Node {
	return nil
}
func (s *SimulationIteratorService) GetSatellites() []types.Satellite {
	return nil
}
func (s *SimulationIteratorService) GetGroundStations() []types.GroundStation {
	return nil
}
func (s *SimulationIteratorService) GetSimulationTime() time.Time {
	return time.Now()
}
func (s *SimulationIteratorService) GetStatePluginRepository() *types.StatePluginRepository {
	return nil
}
func (s *SimulationIteratorService) Close() {

}
