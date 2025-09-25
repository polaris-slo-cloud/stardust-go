package simulation

import (
	"encoding/gob"
	"os"

	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/links"
	"github.com/keniack/stardustGo/internal/node"
	"github.com/keniack/stardustGo/internal/routing"
	"github.com/keniack/stardustGo/pkg/types"
)

// SimulationStateDeserializer is responsible for deserializing simulation state data.
type SimulationStateDeserializer struct {
	inputFile        string
	computingBuilder computing.ComputingBuilder
	routerBuilder    *routing.RouterBuilder
	simPlugins       []types.SimulationPlugin
}

// NewSimulationStateDeserializer creates a new deserializer instance.
func NewSimulationStateDeserializer(inputFile string, computingBuilder computing.ComputingBuilder, routerBuilder *routing.RouterBuilder, simPlugins []types.SimulationPlugin) *SimulationStateDeserializer {
	return &SimulationStateDeserializer{
		inputFile:        inputFile,
		computingBuilder: computingBuilder,
		routerBuilder:    routerBuilder,
		simPlugins:       simPlugins,
	}
}

// Load reads the serialized data from the input file and returns the reconstructed SimulationMetadata.
func (d *SimulationStateDeserializer) load() (*SimulationMetadata, error) {
	file, err := os.Open(d.inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a decoder
	decoder := gob.NewDecoder(file)

	// Decode the data
	var metadata SimulationMetadata
	if err := decoder.Decode(&metadata); err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (d *SimulationStateDeserializer) LoadIterator() types.SimulationController {
	metadata, _ := d.load()
	satellites := make([]types.Satellite, len(metadata.Satellites))
	for i, sat := range metadata.Satellites {
		router, _ := d.routerBuilder.Build()
		computing := d.computingBuilder.Build()
		satellites[i] = node.NewSimulatedSatellite(sat.Name, router, computing, links.NewSimulatedLinkProtocol())
	}

	return nil
}
