package ground

import (
	"log"
	"os"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/links"
	"github.com/keniack/stardustGo/pkg/types"
	"gopkg.in/yaml.v3"
)

type rawGroundStation struct {
	Name          string  `yaml:"Name"`
	Lat           float64 `yaml:"Lat"`
	Lon           float64 `yaml:"Lon"`
	Protocol      string  `yaml:"Protocol"`
	Router        string  `yaml:"Router"`
	ComputingType string  `yaml:"ComputingType"`
}

type GroundStationYmlLoader struct {
	config               configs.GroundLinkConfig
	groundStationBuilder *GroundStationBuilder
}

func NewGroundStationYmlLoader(
	config configs.GroundLinkConfig,
	builder *GroundStationBuilder,
) *GroundStationYmlLoader {
	return &GroundStationYmlLoader{
		config:               config,
		groundStationBuilder: builder,
	}
}

func (l *GroundStationYmlLoader) Load(path string, satellites []types.Satellite) ([]types.GroundStation, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open ground station file: %v", err)
		return nil, err
	}
	defer file.Close()

	var groundStations []rawGroundStation
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&groundStations); err != nil {
		log.Fatalf("Failed to decode ground station YAML: %v", err)
		return nil, err
	}

	var result []types.GroundStation
	for _, gs := range groundStations {
		station := l.groundStationBuilder.
			SetName(gs.Name).
			SetLatitude(gs.Lat).
			SetLongitude(gs.Lon).
			SetComputingType(gs.ComputingType).
			ConfigureGroundLinkProtocol(func(p *links.GroundProtocolBuilder) *links.GroundProtocolBuilder {
				return p.
					SetProtocol(gs.Protocol).
					SetSatellites(satellites)
			}).
			Build()
		result = append(result, station)
	}

	return result, nil
}
