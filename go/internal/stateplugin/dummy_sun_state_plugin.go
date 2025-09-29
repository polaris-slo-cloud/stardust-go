package stateplugin

import (
	"encoding/gob"
	"math/rand"
	"os"
	"reflect"

	"github.com/keniack/stardustGo/pkg/helper"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.StatePlugin = (*DummySunStatePlugin)(nil)
var _ SunStatePlugin = (*DummySunStatePlugin)(nil)

type SunStatePlugin interface {
	types.StatePlugin

	// GetSunlightExposure returns the current sunlight exposure for a satellite (0.0 to 1.0)
	GetSunlightExposure(node types.Node) float64
}

type DummySunStatePlugin struct {
	sunlightExposure map[types.Node]float64 // Map to store satellite -> sunlight exposure (0.0 to 1.0)
	states           []map[string]float64   // array to store the states
}

func NewDummySunStatePlugin() *DummySunStatePlugin {
	return &DummySunStatePlugin{
		sunlightExposure: make(map[types.Node]float64),
	}
}

func (d *DummySunStatePlugin) GetSunlightExposure(node types.Node) float64 {
	return d.sunlightExposure[node]
}

func (d *DummySunStatePlugin) GetName() string {
	return "DummyPlugin"
}

func (d *DummySunStatePlugin) GetType() reflect.Type {
	var dummy SunStatePlugin
	return reflect.TypeOf(dummy)
}

// PostSimulationStep updates the sunlight exposure for each satellite
func (d *DummySunStatePlugin) PostSimulationStep(simulationController types.SimulationController) {
	// ONLY RANDOM DATA FOR DEMO PURPOSES !!!

	// Get all satellites from the simulation
	nodes := simulationController.GetSatellites()

	// For each node, set a random sunlight exposure between 0.0 and 1.0
	for _, node := range nodes {
		// Store random sunlight exposure (0.0 to 1.0) for each satellite
		d.sunlightExposure[node] = rand.Float64()
	}
}

func (d *DummySunStatePlugin) AddState(simulationController types.SimulationController) {
	stateMap := make(map[string]float64)
	for node, state := range d.sunlightExposure {
		stateMap[node.GetName()] = state
	}

	d.states = append(d.states, stateMap)
}

func (d *DummySunStatePlugin) Save(origFile string) {
	filename := helper.ExtendFilename(origFile, ".dummySimPlugin")

	file, _ := os.Create(filename)
	defer file.Close()

	// Create an encoder
	encoder := gob.NewEncoder(file)

	// Encode the data
	encoder.Encode(d.states)
}
