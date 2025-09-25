package simulation

import (
	"encoding/gob"
	"os"

	"github.com/keniack/stardustGo/pkg/types"
)

type SimulationStateSerializer struct {
	outputFile string
	metadata   SimulationMetadata
	linksIxMap map[types.Link]int
}

func NewSimulationStateSerializer(outputFile string) *SimulationStateSerializer {
	return &SimulationStateSerializer{
		outputFile: outputFile,
		metadata:   NewSimulationMetadata(),
		linksIxMap: make(map[types.Link]int),
	}
}

func (s *SimulationStateSerializer) AddState(simulationController types.SimulationController) {
	var nodes = simulationController.GetAllNodes()
	var nodeStates = []NodeState{}
	for _, node := range nodes {
		linkIxs := []int{}
		for _, link := range node.GetLinkNodeProtocol().Established() {
			var linkIx int
			if ix, exists := s.linksIxMap[link]; exists {
				linkIx = ix
			} else {
				n1, n2 := link.Nodes()
				linkIx = len(s.metadata.Links)
				s.linksIxMap[link] = linkIx
				s.metadata.Links = append(s.metadata.Links, SimulationLink{
					NodeName1: n1.GetName(),
					NodeName2: n2.GetName(),
				})
			}
			linkIxs = append(linkIxs, linkIx)
		}
		nodeStates = append(nodeStates, NewNodeState(node.GetName(), node.PositionVector(), linkIxs))
	}
	s.metadata.States = append(s.metadata.States, NewSimulationState(simulationController.GetSimulationTime(), nodeStates))
}

func (s *SimulationStateSerializer) Save() {
	file, _ := os.Create(s.outputFile)
	defer file.Close()

	// Create an encoder
	encoder := gob.NewEncoder(file)

	// Encode the data
	encoder.Encode(s.metadata)
}
