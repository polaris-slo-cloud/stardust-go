package deployment

import (
	"errors"
	"fmt"
)

// DeployableService represents a deployable service with CPU and memory requirements.
type DeployableService struct {
	ServiceName string  // The name of the service
	Cpu         float64 // CPU required by the service
	Memory      float64 // Memory required by the service
}

// NewDeployableService creates a new instance of DeployableService with the specified parameters.
func NewDeployableService(serviceName string, cpu, memory float64) (*DeployableService, error) {
	// Validate the input parameters
	if serviceName == "" {
		return nil, errors.New("serviceName cannot be null or empty")
	}
	if cpu <= 0 {
		return nil, fmt.Errorf("cpu must be greater than zero, got %f", cpu)
	}
	if memory <= 0 {
		return nil, fmt.Errorf("memory must be greater than zero, got %f", memory)
	}

	// Create and return the DeployableService instance
	return &DeployableService{
		ServiceName: serviceName,
		Cpu:         cpu,
		Memory:      memory,
	}, nil
}
