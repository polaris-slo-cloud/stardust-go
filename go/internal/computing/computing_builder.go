package computing

import (
	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ ComputingBuilder = (*DefaultComputingBuilder)(nil)

// ComputingBuilder is the interface for building Computing instances.
type ComputingBuilder interface {
	// WithComputingType sets the computing type for the builder.
	WithComputingType(computingType types.ComputingType) ComputingBuilder

	// Build creates and returns the final Computing instance.
	Build() *Computing // Return a pointer to Computing
}

// DefaultComputingBuilder builds a Computing instance based on a given configuration.
type DefaultComputingBuilder struct {
	computingConfiguration []configs.ComputingConfig
	currentConfiguration   configs.ComputingConfig
}

// NewComputingBuilder creates a new instance of ComputingBuilder with the given configuration.
func NewComputingBuilder(computingConfiguration []configs.ComputingConfig) *DefaultComputingBuilder {
	return &DefaultComputingBuilder{
		computingConfiguration: computingConfiguration,
		currentConfiguration:   computingConfiguration[0], // No computing selected initially
	}
}

// WithComputingType configures the Computing instance with a specific ComputingType.
func (b *DefaultComputingBuilder) WithComputingType(computingType types.ComputingType) ComputingBuilder {
	if b.currentConfiguration.Type == computingType {
		return b
	}
	for _, config := range b.computingConfiguration {
		if config.Type == computingType {
			b.currentConfiguration = config
			break
		}
	}
	return b
}

// Build returns the configured Computing instance.
func (b *DefaultComputingBuilder) Build() *Computing {
	return NewComputing(float64(b.currentConfiguration.Cores), float64(b.currentConfiguration.Memory), b.currentConfiguration.Type)
}
