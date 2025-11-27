package simulation

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"os"

	"github.com/keniack/stardustGo/pkg/types"
)

// SimulationStateSerializer is responsible for serializing the state of a simulation.
// (including StatePlugins)
type SimulationStateSerializer struct {
	outputFile   string
	metadata     types.SimulationMetadata
	linksIxMap   map[types.Link]int
	statePlugins []types.StatePlugin
}

// NewSimulationStateSerializer initializes a new SimulationStateSerializer.
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
	satellites := simualtionController.GetSatellites()
	s.metadata.Satellites = make([]types.RawSatellite, len(satellites))
	for i, sat := range satellites {
		s.metadata.Satellites[i] = types.RawSatellite{
			Index:         i,
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

	s.metadata.StatePlugins = s.metadata.StatePlugins[:0]
	for _, plugin := range s.statePlugins {
		s.metadata.StatePlugins = append(s.metadata.StatePlugins, plugin.GetName())
	}

	// gob output
	file, err := os.Create(s.outputFile)
	if err != nil {
		log.Printf("error creating gob file %s: %v", s.outputFile, err)
		return
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(s.metadata); err != nil {
		log.Printf("error encoding gob state: %v", err)
	}

	// JSON output, same base name plus .json
	jsonPath := s.outputFile + ".json"
	jsonFile, err := os.Create(jsonPath)
	if err != nil {
		log.Printf("error creating json file %s: %v", jsonPath, err)
	} else {
		defer jsonFile.Close()
		jsonEnc := json.NewEncoder(jsonFile)
		jsonEnc.SetIndent("", "  ")
		if err := jsonEnc.Encode(s.metadata); err != nil {
			log.Printf("error encoding json state: %v", err)
		}
	}

	for _, plugin := range s.statePlugins {
		plugin.Save(s.outputFile)
	}
}
