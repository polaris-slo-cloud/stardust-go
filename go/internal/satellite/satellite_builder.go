package satellite

import (
	"time"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/links"
	"github.com/keniack/stardustGo/internal/node"
	"github.com/keniack/stardustGo/internal/routing"
)

// SatelliteBuilder helps construct Satellite instances with ISL, routing, and computing configuration.
type SatelliteBuilder struct {
	name              string
	inclination       float64
	rightAscension    float64
	eccentricity      float64
	argumentOfPerigee float64
	meanAnomaly       float64
	meanMotion        float64
	epoch             time.Time

	routerBuilder    *routing.RouterBuilder
	computingBuilder *computing.DefaultComputingBuilder
	islBuilder       *links.IslProtocolBuilder
	islConfig        configs.InterSatelliteLinkConfig // Store the ISL config
}

// NewSatelliteBuilder creates a new SatelliteBuilder with required dependencies.
func NewSatelliteBuilder(router *routing.RouterBuilder, computing *computing.DefaultComputingBuilder, islConfig configs.InterSatelliteLinkConfig) *SatelliteBuilder {
	return &SatelliteBuilder{
		routerBuilder:    router,
		computingBuilder: computing,
		islConfig:        islConfig, // Initialize the ISL config
	}
}

func (b *SatelliteBuilder) SetName(name string) *SatelliteBuilder {
	b.name = name
	return b
}

func (b *SatelliteBuilder) SetInclination(value float64) *SatelliteBuilder {
	b.inclination = value
	return b
}

func (b *SatelliteBuilder) SetRightAscension(value float64) *SatelliteBuilder {
	b.rightAscension = value
	return b
}

func (b *SatelliteBuilder) SetEccentricity(value float64) *SatelliteBuilder {
	b.eccentricity = value
	return b
}

func (b *SatelliteBuilder) SetArgumentOfPerigee(value float64) *SatelliteBuilder {
	b.argumentOfPerigee = value
	return b
}

func (b *SatelliteBuilder) SetMeanAnomaly(value float64) *SatelliteBuilder {
	b.meanAnomaly = value
	return b
}

func (b *SatelliteBuilder) SetMeanMotion(value float64) *SatelliteBuilder {
	b.meanMotion = value
	return b
}

func (b *SatelliteBuilder) SetEpoch(epoch time.Time) *SatelliteBuilder {
	b.epoch = epoch
	return b
}

// ConfigureISL now uses the ISL config passed to the builder
func (b *SatelliteBuilder) ConfigureISL(fn func(builder *links.IslProtocolBuilder) *links.IslProtocolBuilder) *SatelliteBuilder {
	// Pass the ISL config to the builder
	b.islBuilder = fn(links.NewIslProtocolBuilder(b.islConfig))
	return b
}

// Build constructs the Satellite instance from configured parameters.
func (b *SatelliteBuilder) Build() *node.Satellite {
	// Handle the error returned by routerBuilder.Build()
	router, err := b.routerBuilder.Build()
	if err != nil {
		// Handle error (e.g., log it or return nil)
		// For now, we'll panic, but you can return an error or a fallback value
		panic("failed to build router: " + err.Error())
	}

	return node.NewSatellite(
		b.name,
		b.inclination,
		b.rightAscension,
		b.eccentricity,
		b.argumentOfPerigee,
		b.meanAnomaly,
		b.meanMotion,
		b.epoch,
		time.Now(), // Simulated current time, adjust if needed
		b.islBuilder.Build(),
		router, // Pass the router after error handling
		b.computingBuilder.WithComputingType(configs.ComputingType(configs.Edge)).Build(),
	)
}
