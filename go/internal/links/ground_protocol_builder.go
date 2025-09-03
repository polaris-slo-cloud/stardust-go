package links

import (
	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
)

type GroundProtocolBuilder struct {
	config     configs.GroundLinkConfig
	satellites []types.Node
}

func NewGroundProtocolBuilder(config configs.GroundLinkConfig, satellites []types.Node) *GroundProtocolBuilder {
	return &GroundProtocolBuilder{
		config:     config,
		satellites: satellites,
	}
}

func (b *GroundProtocolBuilder) SetProtocol(protocol string) *GroundProtocolBuilder {
	b.config.Protocol = protocol
	return b
}

func (b *GroundProtocolBuilder) Build() types.GroundSatelliteLinkProtocol {
	switch b.config.Protocol {
	case "nearest":
		return NewGroundSatelliteNearestProtocol(b.satellites)
	default:
		return nil
	}
}
