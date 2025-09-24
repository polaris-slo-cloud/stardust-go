package satellite

import (
	"log"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/simulation"
	"github.com/keniack/stardustGo/pkg/types"
)

// SatelliteLoaderService wires the constellation loader and triggers simulation startup.
type SatelliteLoaderService struct {
	controller            simulation.SimulationController
	constellationLoader   *SatelliteConstellationLoader
	tleLoader             *TleLoader
	satelliteBuilder      *SatelliteBuilder
	config                configs.InterSatelliteLinkConfig
	satelliteDataSource   string
	satelliteSourceFormat string
}

// NewSatelliteLoaderService initializes all required loaders and binds them.
func NewSatelliteLoaderService(
	config configs.InterSatelliteLinkConfig,
	builder *SatelliteBuilder,
	loader *SatelliteConstellationLoader,
	controller simulation.SimulationController,
	dataSourcePath string,
	sourceFormat string,
) *SatelliteLoaderService {
	tleLoader := NewTleLoader(config, builder)
	loader.RegisterDataSourceLoader("tle", tleLoader)

	return &SatelliteLoaderService{
		controller:            controller,
		constellationLoader:   loader,
		tleLoader:             tleLoader,
		satelliteBuilder:      builder,
		config:                config,
		satelliteDataSource:   dataSourcePath,
		satelliteSourceFormat: sourceFormat,
	}
}

// Start loads satellites and injects them into the simulation
func (s *SatelliteLoaderService) Start() error {
	log.Println("Starting LoaderService...")
	satellites, err := s.constellationLoader.LoadSatelliteConstellation(s.satelliteDataSource, s.satelliteSourceFormat)
	if err != nil {
		return err
	}

	// Convert []*node.Satellite to []*types.Node
	var nodes []types.Node
	for _, satellite := range satellites {
		// Append the pointer to the slice
		node := types.Node(satellite) // Convert *node.Satellite to *types.Node
		nodes = append(nodes, node)   // Append pointer to slice
	}

	log.Printf("Injecting %d satellites into simulation", len(satellites))
	return s.controller.InjectSatellites(nodes)
}
