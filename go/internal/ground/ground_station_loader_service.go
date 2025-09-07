package ground

import (
	"log"

	"github.com/keniack/stardustGo/internal/simulation"
	"github.com/keniack/stardustGo/pkg/types"
)

type GroundStationLoaderService struct {
	controller                simulation.SimulationController
	groundStationBuilder      *GroundStationBuilder
	groundStationLoader       *GroundStationYmlLoader
	groundStationDataSource   string
	groundStationSourceFormat string
}

func NewGroundStationLoaderService(
	controller simulation.SimulationController,
	builder *GroundStationBuilder,
	groundStationLoader *GroundStationYmlLoader,
	dataSourcePath string,
	sourceFormat string,
) *GroundStationLoaderService {
	return &GroundStationLoaderService{
		controller:                controller,
		groundStationBuilder:      builder,
		groundStationLoader:       groundStationLoader,
		groundStationDataSource:   dataSourcePath,
		groundStationSourceFormat: sourceFormat,
	}
}

func (s *GroundStationLoaderService) Start() error {
	log.Println("Starting LoaderService...")
	groundStations, err := s.groundStationLoader.Load(s.groundStationDataSource)
	if err != nil {
		return err
	}

	// Convert []*node.Satellite to []*types.Node
	var nodes []types.Node
	for _, gs := range groundStations {
		// Append the pointer to the slice
		node := types.Node(gs)      // Convert *node.GroundStation to *types.Node
		nodes = append(nodes, node) // Append pointer to slice
	}

	log.Printf("Injecting %d ground stations into simulation", len(groundStations))
	return s.controller.InjectGroundStations(nodes)
}
