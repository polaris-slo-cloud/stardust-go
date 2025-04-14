package satellite

import (
	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
	"log"
)

// LoaderService wires the constellation loader and triggers simulation startup.
type LoaderService struct {
	controller            types.ISimulationController
	constellationLoader   *SatelliteConstellationLoader
	tleLoader             *TleLoader
	satelliteBuilder      *SatelliteBuilder
	config                configs.InterSatelliteLinkConfig
	satelliteDataSource   string
	satelliteSourceFormat string
}

// NewLoaderService initializes all required loaders and binds them.
func NewLoaderService(
	config configs.InterSatelliteLinkConfig,
	builder *SatelliteBuilder,
	loader *SatelliteConstellationLoader,
	controller types.ISimulationController,
	dataSourcePath string,
	sourceFormat string,
) *LoaderService {
	tleLoader := NewTleLoader(config, builder)
	loader.RegisterDataSourceLoader("tle", tleLoader)

	return &LoaderService{
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
func (s *LoaderService) Start() error {
	log.Println("Starting LoaderService...")
	satellites, err := s.constellationLoader.LoadSatelliteConstellation(s.satelliteDataSource, s.satelliteSourceFormat)
	if err != nil {
		return err
	}

	// Convert []*node.Satellite to []*types.INode
	var nodes []types.INode
	for _, satellite := range satellites {
		// Append the pointer to the slice
		node := types.INode(satellite) // Convert *node.Satellite to *types.INode
		nodes = append(nodes, node)    // Append pointer to slice
	}

	log.Printf("Injecting %d satellites into simulation", len(satellites))
	return s.controller.InjectSatellites(nodes)
}
