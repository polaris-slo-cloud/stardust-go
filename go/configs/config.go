package configs

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/keniack/stardustGo/pkg/types"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Simulation SimulationConfig         `json:"SimulationConfiguration" yaml:"SimulationConfiguration"`
	ISL        InterSatelliteLinkConfig `json:"InterSatelliteLinkConfig" yaml:"InterSatelliteLinkConfig"`
	Ground     GroundLinkConfig         `json:"GroundLinkConfig" yaml:"GroundLinkConfig"`
	Router     RouterConfig             `json:"RouterConfig" yaml:"RouterConfig"`
	Computing  []ComputingConfig        `json:"ComputingConfiguration" yaml:"ComputingConfiguration"`
}

type SimulationConfig struct {
	StepInterval                int       `json:"StepInterval" yaml:"StepInterval"`
	StepMultiplier              int       `json:"StepMultiplier" yaml:"StepMultiplier"`
	SatelliteDataSource         string    `json:"SatelliteDataSource" yaml:"SatelliteDataSource"`
	SatelliteDataSourceType     string    `json:"SatelliteDataSourceType" yaml:"SatelliteDataSourceType"`
	GroundStationDataSource     string    `json:"GroundStationDataSource" yaml:"GroundStationDataSource"`
	GroundStationDataSourceType string    `json:"GroundStationDataSourceType" yaml:"GroundStationDataSourceType"`
	UsePreRouteCalc             bool      `json:"UsePreRouteCalc" yaml:"UsePreRouteCalc"`
	MaxCpuCores                 int       `json:"MaxCpuCores" yaml:"MaxCpuCores"`
	SimulationStartTime         time.Time `json:"SimulationStartTime" yaml:"SimulationStartTime"`
	Plugins                     []string  `json:"Plugins" yaml:"Plugins"`
}

type InterSatelliteLinkConfig struct {
	Neighbours int    `json:"Neighbours" yaml:"Neighbours"` // Number of neighbors per satellite
	Protocol   string `json:"Protocol" yaml:"Protocol"`     // Strategy name: "mst", "nearest", etc.
}

type GroundLinkConfig struct {
	Protocol string `json:"Protocol" yaml:"Protocol"`
}

type RouterConfig struct {
	Protocol string `json:"Protocol" yaml:"Protocol"`
}

type ComputingConfig struct {
	Cores  int                 `json:"Cores" yaml:"Cores"`
	Memory int                 `json:"Memory" yaml:"Memory"`
	Type   types.ComputingType `json:"Type" yaml:"Type"` // Should be either "Edge" or "Cloud"
}

func LoadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	case ".json":
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported config file type")
	}
	return &cfg, nil
}
