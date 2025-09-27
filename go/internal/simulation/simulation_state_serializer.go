package simulation

import (
	"encoding/gob"
	"os"

	"github.com/keniack/stardustGo/pkg/types"
)

type SimulationStateSerializer struct {
	outputFile    string
	metadata      SimulationMetadata
	linksIxMap    map[types.Link]int
	nodeNodeIxMap map[types.Node]map[types.Node]int
}

func NewSimulationStateSerializer(outputFile string) *SimulationStateSerializer {
	return &SimulationStateSerializer{
		outputFile:    outputFile,
		metadata:      NewSimulationMetadata(),
		linksIxMap:    make(map[types.Link]int),
		nodeNodeIxMap: make(map[types.Node]map[types.Node]int),
	}
}

func (s *SimulationStateSerializer) AddState(simulationController types.SimulationController) {
	var nodes = simulationController.GetAllNodes()
	if len(s.metadata.States) == 0 {
		for _, node := range nodes {
			s.nodeNodeIxMap[node] = make(map[types.Node]int)
		}
	}

	var nodeStates = []NodeState{}
	for _, node := range nodes {
		linkIxs := []int{}
		for _, link := range node.GetLinkNodeProtocol().Established() {
			var linkIx int
			nn, nm := link.Nodes()
			if ix, exists := s.nodeNodeIxMap[nn][nm]; exists {
				linkIx = ix
			} else {
				n1, n2 := link.Nodes()
				linkIx = len(s.metadata.Links)
				s.linksIxMap[link] = linkIx
				s.nodeNodeIxMap[n1][n2] = linkIx
				s.metadata.Links = append(s.metadata.Links, SimulationLink{
					NodeName1: n1.GetName(),
					NodeName2: n2.GetName(),
				})
			}
			linkIxs = append(linkIxs, linkIx)
		}
		nodeStates = append(nodeStates, NewNodeState(node.GetName(), node.GetPosition(), linkIxs))
	}
	s.metadata.States = append(s.metadata.States, NewSimulationState(simulationController.GetSimulationTime(), nodeStates))
}

func (s *SimulationStateSerializer) Save(simualtionController types.SimulationController) {

	var satellites = simualtionController.GetSatellites()
	s.metadata.Satellites = make([]RawSatellite, len(satellites))
	for i, sat := range satellites {
		s.metadata.Satellites[i] = RawSatellite{
			Name:          sat.GetName(),
			ComputingType: sat.GetComputing().GetComputingType(),
		}
	}

	s.metadata.Grounds = make([]RawGroundStation, len(simualtionController.GetGroundStations()))
	for i, gs := range simualtionController.GetGroundStations() {
		s.metadata.Grounds[i] = RawGroundStation{
			Name:          gs.GetName(),
			ComputingType: gs.GetComputing().GetComputingType(),
		}
	}

	file, _ := os.Create(s.outputFile)
	defer file.Close()

	// Create an encoder
	encoder := gob.NewEncoder(file)

	// Encode the data
	encoder.Encode(s.metadata)
}
