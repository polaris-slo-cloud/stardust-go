package simulation

import (
	"encoding/gob"
	"log"
	"os"
	"time"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/deployment"
	"github.com/keniack/stardustGo/internal/links"
	"github.com/keniack/stardustGo/internal/links/linktypes"
	"github.com/keniack/stardustGo/internal/node"
	"github.com/keniack/stardustGo/internal/routing"
	"github.com/keniack/stardustGo/pkg/types"
)

// SimulationStateDeserializer is responsible for deserializing simulation state data.
type SimulationStateDeserializer struct {
	inputFile        string
	computingBuilder computing.ComputingBuilder
	routerBuilder    *routing.RouterBuilder
	orchestrator     *deployment.DeploymentOrchestrator
	simPlugins       []types.SimulationPlugin
	config           *configs.SimulationConfig
}

// NewSimulationStateDeserializer creates a new deserializer instance.
func NewSimulationStateDeserializer(config *configs.SimulationConfig, inputFile string, computingBuilder computing.ComputingBuilder, routerBuilder *routing.RouterBuilder, orchestrator *deployment.DeploymentOrchestrator, simPlugins []types.SimulationPlugin) *SimulationStateDeserializer {
	return &SimulationStateDeserializer{
		inputFile:        inputFile,
		computingBuilder: computingBuilder,
		routerBuilder:    routerBuilder,
		orchestrator:     orchestrator,
		simPlugins:       simPlugins,
		config:           config,
	}
}

// Load reads the serialized data from the input file and returns the reconstructed SimulationMetadata.
func (d *SimulationStateDeserializer) load() *SimulationMetadata {
	file, err := os.Open(d.inputFile)
	if err != nil {
		log.Fatalln("Failed to open simulation state file:", err)
	}
	defer file.Close()

	// Create a decoder
	decoder := gob.NewDecoder(file)

	// Decode the data
	var metadata SimulationMetadata
	if err := decoder.Decode(&metadata); err != nil {
		log.Fatalln("Failed to decode simulation state:", err)
	}

	return &metadata
}

func (d *SimulationStateDeserializer) LoadIterator() types.SimulationController {
	metadata := d.load()

	innerProtocol := links.NewSimulatedLinkProtocol()

	// Reconstruct nodes
	nodeNames := make(map[string]node.SimulatedNode)
	satellites := make([]types.Node, len(metadata.Satellites))
	for i, sat := range metadata.Satellites {
		router, _ := d.routerBuilder.Build()
		computing := d.computingBuilder.Build()
		satellite := node.NewSimulatedSatellite(sat.Name, router, computing, links.NewLinkFilterProtocol(innerProtocol))
		satellites[i] = satellite
		nodeNames[sat.Name] = satellite
	}

	groundStations := make([]types.Node, len(metadata.Grounds))
	for i, gs := range metadata.Grounds {
		router, _ := d.routerBuilder.Build()
		computing := d.computingBuilder.Build()
		groundStation := node.NewSimulatedGroundStation(gs.Name, router, computing, links.NewLinkFilterProtocol(innerProtocol))
		groundStations[i] = groundStation
		nodeNames[gs.Name] = groundStation
	}

	// Reconstruct links
	links := make([]types.Link, len(metadata.Links))
	for i, l := range metadata.Links {
		var n1, n2 node.SimulatedNode
		n1 = nodeNames[l.NodeName1]
		n2 = nodeNames[l.NodeName2]
		links[i] = linktypes.NewSimulatedLink(n1, n2)
		innerProtocol.AddLink(links[i])
	}

	type positionState struct {
		time     time.Time
		position types.Vector
	}

	// Reconstruct states
	positions := make(map[string][]positionState)
	establishedLinks := make([][]types.Link, len(metadata.States))
	for i, state := range metadata.States {
		linkIxSeen := make(map[int]bool)
		for _, nodeState := range state.NodeStates {
			for _, linkIx := range nodeState.Established {
				if linkIxSeen[linkIx] {
					continue
				}
				establishedLinks[i] = append(establishedLinks[i], links[linkIx])
				linkIxSeen[linkIx] = true
			}
			positions[nodeState.Name] = append(positions[nodeState.Name], positionState{
				time:     state.Time,
				position: nodeState.Position,
			})
		}
	}

	innerProtocol.InjectEstablishedLinks(establishedLinks)
	for name, node := range nodeNames {
		for _, ps := range positions[name] {
			node.AddPositionState(ps.time, ps.position)
		}
	}

	simService := NewSimulationIteratorService(d.config, metadata.States, d.simPlugins)
	simService.Inject(d.orchestrator)
	simService.InjectSatellites(satellites)
	simService.InjectGroundStations(groundStations)

	return simService
}
