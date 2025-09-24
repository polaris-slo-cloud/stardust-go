package ground

import (
	"time"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/computing"
	"github.com/keniack/stardustGo/internal/links"
	"github.com/keniack/stardustGo/internal/routing"
	"github.com/keniack/stardustGo/pkg/types"
)

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

func NewGroundStationBuilder(simStartTime time.Time, router *routing.RouterBuilder, computing *computing.DefaultComputingBuilder, config configs.GroundLinkConfig) *GroundStationBuilder {
	return &GroundStationBuilder{
		simStartTime:     simStartTime,
		routerBuilder:    router,
		computingBuilder: computing,
		protocolBuilder:  links.NewGroundProtocolBuilder(config),
	}
}

func (b *GroundStationBuilder) SetName(name string) *GroundStationBuilder {
	b.name = name
	return b
}

func (b *GroundStationBuilder) SetLatitude(value float64) *GroundStationBuilder {
	b.latitude = value
	return b
}

func (b *GroundStationBuilder) SetLongitude(value float64) *GroundStationBuilder {
	b.longitude = value
	return b
}

func (b *GroundStationBuilder) SetAltitude(value float64) *GroundStationBuilder {
	b.altitude = value
	return b
}

func (b *GroundStationBuilder) SetComputingType(value string) *GroundStationBuilder {
	ctype, _ := configs.ToComputingType(value)
	b.computingBuilder.WithComputingType(ctype)
	return b
}

func (b *GroundStationBuilder) ConfigureGroundLinkProtocol(fn func(*links.GroundProtocolBuilder) *links.GroundProtocolBuilder) *GroundStationBuilder {
	b.protocolBuilder = fn(b.protocolBuilder)
	return b
}

func (b *GroundStationBuilder) Build() *types.GroundStation {
	router, err := b.routerBuilder.Build()
	if err != nil {
		panic(err)
	}

	return types.NewGroundStation(
		b.name,
		b.latitude,
		b.longitude,
		b.protocolBuilder.Build(),
		b.simStartTime,
		router,
		b.computingBuilder.Build())
}
