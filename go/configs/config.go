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

type SimulationConfig struct {
	StepInterval                int       `json:"StepInterval" yaml:"StepInterval"`
	StepMultiplier              int       `json:"StepMultiplier" yaml:"StepMultiplier"`
	StepCount                   int       `json:"StepCount" yaml:"StepCount"`
	SatelliteDataSource         string    `json:"SatelliteDataSource" yaml:"SatelliteDataSource"`
	SatelliteDataSourceType     string    `json:"SatelliteDataSourceType" yaml:"SatelliteDataSourceType"`
	GroundStationDataSource     string    `json:"GroundStationDataSource" yaml:"GroundStationDataSource"`
	GroundStationDataSourceType string    `json:"GroundStationDataSourceType" yaml:"GroundStationDataSourceType"`
	UsePreRouteCalc             bool      `json:"UsePreRouteCalc" yaml:"UsePreRouteCalc"`
	SimulationStartTime         time.Time `json:"SimulationStartTime" yaml:"SimulationStartTime"`
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

// LoadConfigFromFile loads a configuration of type T from a file.
// Supported file types: .yaml, .yml, .json
func LoadConfigFromFile[T any](path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg T
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
