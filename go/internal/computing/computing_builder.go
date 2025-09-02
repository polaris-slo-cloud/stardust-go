package computing

import (
	"github.com/keniack/stardustGo/configs"
)

var _ ComputingBuilder = (*DefaultComputingBuilder)(nil)

// ComputingBuilder is the interface for building Computing instances.
type ComputingBuilder interface {
	// WithComputingType sets the computing type for the builder.
	WithComputingType(computingType configs.ComputingType) ComputingBuilder

	// Build creates and returns the final Computing instance.
	Build() *Computing // Return a pointer to Computing
}

// DefaultComputingBuilder builds a Computing instance based on a given configuration.
type DefaultComputingBuilder struct {
	computingConfiguration configs.ComputingConfig
	useComputing           *Computing // Pointer to Computing
}

// NewComputingBuilder creates a new instance of ComputingBuilder with the given configuration.
func NewComputingBuilder(computingConfiguration configs.ComputingConfig) *DefaultComputingBuilder {
	return &DefaultComputingBuilder{
		computingConfiguration: computingConfiguration,
		useComputing:           nil, // No computing selected initially
	}
}

// WithComputingType configures the Computing instance with a specific ComputingType.
func (b *DefaultComputingBuilder) WithComputingType(computingType configs.ComputingType) ComputingBuilder {
	// Set the computing type in the computing configuration
	b.computingConfiguration.Type = computingType

	// Based on the selected computing type, we configure the computing unit.
	switch computingType {
	case configs.Edge:
		// Set parameters for Edge computing
		b.useComputing = NewComputing(8, 16, computingType) // Example values, adjust as needed
	case configs.Cloud:
		// Set parameters for Cloud computing
		b.useComputing = NewComputing(16, 32, computingType) // Example values, adjust as needed
	default:
		// Set default computing parameters
		b.useComputing = NewComputing(4, 8, computingType) // Default example values
	}

	return b
}

// Build returns the configured Computing instance.
func (b *DefaultComputingBuilder) Build() *Computing {
	return b.useComputing // Return the pointer to Computing
}
