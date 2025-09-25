package simulation

import (
	"time"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.SimulationController = (*SimulationIteratorService)(nil)

type SimulationIteratorService struct {
	BaseSimulationService
}

func NewSimulationIteratorService(config *configs.SimulationConfig) *SimulationIteratorService {
	service := &SimulationIteratorService{}
	service.BaseSimulationService = NewBaseSimulationService(config, service.runSimulationStep)
	return service
}

func (s *SimulationIteratorService) GetStatePluginRepository() *types.StatePluginRepository {
	return nil
}

func (s *SimulationIteratorService) Close() {

}

func (s *SimulationIteratorService) runSimulationStep(nextTime func(time.Time) time.Time) {

}
