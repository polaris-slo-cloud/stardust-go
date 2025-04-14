package configs

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Simulation SimulationConfig         `json:"SimulationConfiguration"`
	ISL        InterSatelliteLinkConfig `json:"InterSatelliteLinkConfig"`
	Router     RouterConfig             `json:"RouterConfig"`
	Computing  []ComputingConfig        `json:"ComputingConfiguration"`
}

type SimulationConfig struct {
	StepInterval            int       `json:"StepInterval"`
	StepMultiplier          int       `json:"StepMultiplier"`
	SatelliteDataSource     string    `json:"SatelliteDataSource"`
	SatelliteDataSourceType string    `json:"SatelliteDataSourceType"`
	UsePreRouteCalc         bool      `json:"UsePreRouteCalc"`
	MaxCpuCores             int       `json:"MaxCpuCores"`
	SimulationStartTime     time.Time `json:"SimulationStartTime"`
}

type InterSatelliteLinkConfig struct {
	Neighbours int    `json:"neighbours"` // Number of neighbors per satellite
	Protocol   string `json:"protocol"`   // Strategy name: "mst", "nearest", etc.
}

type RouterConfig struct {
	Protocol string `json:"protocol" yaml:"protocol"`
}

type ComputingConfig struct {
	Cores  int           `json:"Cores"`
	Memory int           `json:"Memory"`
	Type   ComputingType `json:"Type"` // Should be either "Edge" or "Cloud"
}

type ComputingConfigList struct {
	Type []ComputingConfig // This should hold multiple computing configurations
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

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
