package ground

import (
	"time"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/links"
	"github.com/keniack/stardustGo/internal/node"
	"github.com/keniack/stardustGo/internal/routing"
	"github.com/keniack/stardustGo/pkg/types"
)

// GroundStationBuilder is a builder pattern implementation for creating ground stations.
type GroundStationBuilder struct {
	name      string
	latitude  float64
	longitude float64
	altitude  float64

	simStartTime     time.Time
	protocolBuilder  *links.GroundProtocolBuilder
	routerBuilder    *routing.RouterBuilder
	computingBuilder *computing.DefaultComputingBuilder
}

// NewGroundStationBuilder initializes a new GroundStationBuilder
func NewGroundStationBuilder(simStartTime time.Time, router *routing.RouterBuilder, computing *computing.DefaultComputingBuilder, config configs.GroundLinkConfig) *GroundStationBuilder {
	return &GroundStationBuilder{
		simStartTime:     simStartTime,
		routerBuilder:    router,
		computingBuilder: computing,
		protocolBuilder:  links.NewGroundProtocolBuilder(config),
	}
}

// SetName sets the name of the ground station and returns the builder for chaining.
func (b *GroundStationBuilder) SetName(name string) *GroundStationBuilder {
	b.name = name
	return b
}

// SetLatitude sets the latitude coordinate of the ground station and returns the builder for chaining.
func (b *GroundStationBuilder) SetLatitude(value float64) *GroundStationBuilder {
	b.latitude = value
	return b
}

// SetLongitude sets the longitude coordinate of the ground station and returns the builder for chaining.
func (b *GroundStationBuilder) SetLongitude(value float64) *GroundStationBuilder {
	b.longitude = value
	return b
}

// SetAltitude sets the altitude of the ground station and returns the builder for chaining.
func (b *GroundStationBuilder) SetAltitude(value float64) *GroundStationBuilder {
	b.altitude = value
	return b
}

// SetComputingType sets the computing type for the ground station and returns the builder for chaining.
func (b *GroundStationBuilder) SetComputingType(value string) *GroundStationBuilder {
	ctype, _ := types.ToComputingType(value)
	b.computingBuilder.WithComputingType(ctype)
	return b
}

// ConfigureGroundLinkProtocol allows for custom configuration of the ground link protocol.
func (b *GroundStationBuilder) ConfigureGroundLinkProtocol(fn func(*links.GroundProtocolBuilder) *links.GroundProtocolBuilder) *GroundStationBuilder {
	b.protocolBuilder = fn(b.protocolBuilder)
	return b
}

// Build constructs and returns a new GroundStation using the configured properties.
// It panics if the router cannot be built.
func (b *GroundStationBuilder) Build() types.GroundStation {
	router, err := b.routerBuilder.Build()
	if err != nil {
		panic(err)
	}

	return node.NewGroundStation(
		b.name,
		b.latitude,
		b.longitude,
		b.protocolBuilder.Build(),
		b.simStartTime,
		router,
		b.computingBuilder.Build())
}
