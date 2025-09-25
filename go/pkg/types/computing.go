package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Computing is an interface for computing resources, including managing services, available resources, and service placement.
type Computing interface {
	// Mount attaches the computing unit to a node
	Mount(node *Node) error

	// GetComputingType return the type of computing resource
	GetComputingType() ComputingType

	// TryPlaceDeploymentAsync tries to place a service on this computing unit
	TryPlaceDeploymentAsync(service DeployableService) (bool, error)

	// RemoveDeploymentAsync removes a deployed service from the computing unit
	RemoveDeploymentAsync(service DeployableService) error

	// CanPlace checks if the service can be placed on this computing unit
	CanPlace(service DeployableService) bool

	// HostsService checks if the computing unit hosts a service by name
	HostsService(serviceName string) bool

	// CpuAvailable returns the remaining CPU available
	CpuAvailable() float64

	// MemoryAvailable returns the remaining memory available
	MemoryAvailable() float64

	// Clone creates a new copy of the current computing unit
	Clone() Computing

	// GetServices returns the list of deployed services
	GetServices() []DeployableService
}

// ComputingType represents the type of computing resource.
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

// ToComputingType converts a string to a ComputingType.
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
