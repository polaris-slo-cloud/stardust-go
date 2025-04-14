package computing

import (
	"stardustGo/configs"
)

// IComputingBuilder is the interface for building Computing instances.
type IComputingBuilder interface {
	// WithComputingType sets the computing type for the builder.
	WithComputingType(computingType configs.ComputingType) IComputingBuilder

	// Build creates and returns the final Computing instance.
	Build() *Computing // Return a pointer to Computing
}

// ComputingBuilder builds a Computing instance based on a given configuration.
type ComputingBuilder struct {
	computingConfiguration configs.ComputingConfig
	useComputing           *Computing // Pointer to Computing
}

// NewComputingBuilder creates a new instance of ComputingBuilder with the given configuration.
func NewComputingBuilder(computingConfiguration configs.ComputingConfig) *ComputingBuilder {
	return &ComputingBuilder{
		computingConfiguration: computingConfiguration,
		useComputing:           nil, // No computing selected initially
	}
}

// WithComputingType configures the Computing instance with a specific ComputingType.
func (b *ComputingBuilder) WithComputingType(computingType configs.ComputingType) IComputingBuilder {
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
func (b *ComputingBuilder) Build() *Computing {
	return b.useComputing // Return the pointer to Computing
}
