package node

import (
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.Satellite = (*SimulatedSatellite)(nil)

type SimulatedSatellite struct {
	BaseNode

	ISLProtocol types.InterSatelliteLinkProtocol
}

func NewSimulatedSatellite(name string, router types.Router, computing types.Computing, isl types.InterSatelliteLinkProtocol) *SimulatedSatellite {
	return &SimulatedSatellite{

		BaseNode:    BaseNode{Name: name, Router: router, Computing: computing},
		ISLProtocol: isl,
	}
}

func (s *SimulatedSatellite) UpdatePosition(simTime time.Time) {
	// no-op for simulation iterator
}

func (s *SimulatedSatellite) GetLinkNodeProtocol() types.LinkNodeProtocol {
	return s.ISLProtocol
}

func (s *SimulatedSatellite) GetISLProtocol() types.InterSatelliteLinkProtocol {
	return s.ISLProtocol
}
