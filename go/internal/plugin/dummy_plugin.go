package plugin

import (
	"log"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.SimulationPlugin = (*DummyPlugin)(nil)

type DummyPlugin struct {
}

func (p *DummyPlugin) Name() string {
	return "DummyPlugin"
}

func (p *DummyPlugin) PostSimulationStep(simulation types.SimulationController) error {
	log.Println("DummyPlugin: PostSimulationStep called")
	log.Println("Current Simulation Time:", simulation.GetSimulationTime())
	log.Println("Number of Nodes:", len(simulation.GetAllNodes()))
	log.Println("Number of Satellites:", len(simulation.GetSatellites()))
	log.Println("Number of Ground Stations:", len(simulation.GetGroundStations()))
	return nil
}
