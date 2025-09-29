package simulation

import (
	"encoding/gob"
	"os"

	"github.com/keniack/stardustGo/pkg/types"
)

type SimulationStateSerializer struct {
	outputFile   string
	metadata     types.SimulationMetadata
	linksIxMap   map[types.Link]int
	statePlugins []types.StatePlugin
}

func NewSimulationStateSerializer(outputFile string, statePlugins []types.StatePlugin) *SimulationStateSerializer {
	return &SimulationStateSerializer{
		outputFile:   outputFile,
		metadata:     types.NewSimulationMetadata(),
		linksIxMap:   make(map[types.Link]int),
		statePlugins: statePlugins,
	}
}

func (s *SimulationStateSerializer) AddState(simulationController types.SimulationController) {
	var nodes = simulationController.GetAllNodes()
	var nodeStates = []types.NodeState{}
	for _, node := range nodes {
		established := node.GetLinkNodeProtocol().Established()
		linkIxs := make([]int, len(established))
		for i, link := range established {
			var linkIx int
			n1, n2 := link.Nodes()
			if ix, exists := s.linksIxMap[link]; exists {
				linkIx = ix
			} else {
				linkIx = len(s.metadata.Links)
				s.linksIxMap[link] = linkIx
				s.metadata.Links = append(s.metadata.Links, types.SimulationLink{
					NodeName1: n1.GetName(),
					NodeName2: n2.GetName(),
				})
			}
			linkIxs[i] = linkIx
		}
		nodeStates = append(nodeStates, types.NewNodeState(node.GetName(), node.GetPosition(), linkIxs))
	}
	s.metadata.States = append(s.metadata.States, types.NewSimulationState(simulationController.GetSimulationTime(), nodeStates))

	for _, plugin := range s.statePlugins {
		plugin.AddState(simulationController)
	}
}

func (s *SimulationStateSerializer) Save(simualtionController types.SimulationController) {

	var satellites = simualtionController.GetSatellites()
	s.metadata.Satellites = make([]types.RawSatellite, len(satellites))
	for i, sat := range satellites {
		s.metadata.Satellites[i] = types.RawSatellite{
			Name:          sat.GetName(),
			ComputingType: sat.GetComputing().GetComputingType(),
		}
	}

	s.metadata.Grounds = make([]types.RawGroundStation, len(simualtionController.GetGroundStations()))
	for i, gs := range simualtionController.GetGroundStations() {
		s.metadata.Grounds[i] = types.RawGroundStation{
			Name:          gs.GetName(),
			ComputingType: gs.GetComputing().GetComputingType(),
		}
	}

	for _, plugin := range s.statePlugins {
		s.metadata.StatePlugins = append(s.metadata.StatePlugins, plugin.GetName())
	}

	file, _ := os.Create(s.outputFile)
	defer file.Close()

	// Create an encoder
	encoder := gob.NewEncoder(file)

	// Encode the data
	encoder.Encode(s.metadata)

	for _, plugin := range s.statePlugins {
		plugin.Save(s.outputFile)
	}
}
