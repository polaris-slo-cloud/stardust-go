package node

import (
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.Satellite = (*PrecomputedSatellite)(nil)

type PrecomputedSatellite struct {
	BaseNode

	ISLProtocol types.InterSatelliteLinkProtocol
	positions   map[time.Time]types.Vector
}

func NewSimulatedSatellite(name string, router types.Router, computing types.Computing, isl types.InterSatelliteLinkProtocol) *PrecomputedSatellite {
	satellite := &PrecomputedSatellite{
		BaseNode:    BaseNode{Name: name, Router: router, Computing: computing},
		ISLProtocol: isl,
		positions:   make(map[time.Time]types.Vector),
	}

	isl.Mount(satellite)
	router.Mount(satellite)
	computing.Mount(satellite)
	return satellite
}

func (s *PrecomputedSatellite) UpdatePosition(simTime time.Time) {
	s.Position = s.positions[simTime]
}

func (s *PrecomputedSatellite) GetLinkNodeProtocol() types.LinkNodeProtocol {
	return s.ISLProtocol
}

func (s *PrecomputedSatellite) GetISLProtocol() types.InterSatelliteLinkProtocol {
	return s.ISLProtocol
}

func (s *PrecomputedSatellite) AddPositionState(time time.Time, position types.Vector) {
	s.positions[time] = position
}
