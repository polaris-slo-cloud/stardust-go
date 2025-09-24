package configs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	Cores  int           `json:"Cores" yaml:"Cores"`
	Memory int           `json:"Memory" yaml:"Memory"`
	Type   ComputingType `json:"Type" yaml:"Type"` // Should be either "Edge" or "Cloud"
}

type ComputingType int

const (
	// None represents an undefined computing type.
	None ComputingType = iota
	// Edge represents edge computing resources.
	Edge
	// Cloud represents cloud computing resources.
	Cloud
	// Any represents any available computing type.
	Any
)

// String converts the ComputingType to a string representation.
func (c ComputingType) String() string {
	return [...]string{"None", "Edge", "Cloud", "Any"}[c]
}

func ToComputingType(s string) (ComputingType, error) {
	switch strings.ToLower(s) {
	case "none":
		return None, nil
	case "edge":
		return Edge, nil
	case "cloud":
		return Cloud, nil
	case "any":
		return Any, nil
	default:
		return None, fmt.Errorf("unknown ComputingType: %s", s)
	}
}

// UnmarshalJSON allows ComputingType to be parsed from JSON as a string.
func (c *ComputingType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	ct, err := ToComputingType(s)
	if err != nil {
		return err
	}

	*c = ct
	return nil
}

// UnmarshalYAML allows ComputingType to be parsed from YAML as a string.
func (c *ComputingType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	ct, err := ToComputingType(s)
	if err != nil {
		return err
	}

	*c = ct
	return nil
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
