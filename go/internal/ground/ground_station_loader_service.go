package ground

import (
	"log"

	"github.com/keniack/stardustGo/pkg/types"
)

// GroundStationLoaderService is responsible for loading ground station configurations
// from a specified data source and injecting them into the simulation controller.
type GroundStationLoaderService struct {
	controller                types.SimulationController
	groundStationBuilder      *GroundStationBuilder
	groundStationLoader       *GroundStationYmlLoader
	groundStationDataSource   string
	groundStationSourceFormat string
}

// NewGroundStationLoaderService initializes a new GroundStationLoaderService.
func NewGroundStationLoaderService(
	controller types.SimulationController,
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

// Start loads ground station configurations from the data source, converts them to Node types,
// and injects them into the simulation controller.
// Returns an error if the loading or injection process fails.
func (s *GroundStationLoaderService) Start() error {
	log.Println("Starting LoaderService...")
	groundStations, err := s.groundStationLoader.Load(s.groundStationDataSource, s.controller.GetSatellites())
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

	return s.controller.InjectGroundStations(nodes)
}
