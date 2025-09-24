package stateplugin

import (
	"math/rand"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.StatePlugin = (*DummySunStatePlugin)(nil)

type DummySunStatePlugin struct {
	sunlightExposure map[types.Node]float64 // Map to store satellite name -> sunlight exposure (0.0 to 1.0)
}

func NewDummySunStatePlugin() *DummySunStatePlugin {
	return &DummySunStatePlugin{
		sunlightExposure: make(map[types.Node]float64),
	}
}

// PostSimulationStep updates the sunlight exposure for each satellite
func (d *DummySunStatePlugin) PostSimulationStep(simulationController types.SimulationController) {
	// ONLY RANDOM DATA FOR DEMO PURPOSES !!!

	// Get all satellites from the simulation
	nodes := simulationController.GetAllNodes()

	// For each node, set a random sunlight exposure between 0.0 and 1.0
	for _, node := range nodes {
		// Store random sunlight exposure (0.0 to 1.0) for each satellite
		d.sunlightExposure[node] = rand.Float64()
	}
}

// GetSunlightExposure returns the current sunlight exposure for a satellite (0.0 to 1.0)
func (d *DummySunStatePlugin) GetSunlightExposure(node types.Node) float64 {
	return d.sunlightExposure[node]
}
