package node

import (
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.GroundStation = (*PrecomputedGroundStation)(nil)

type PrecomputedGroundStation struct {
	BaseNode

	LinkProtocol types.LinkNodeProtocol
	positions    map[time.Time]types.Vector
}

func NewSimulatedGroundStation(name string, router types.Router, computing types.Computing, linkProtocol types.LinkNodeProtocol) *PrecomputedGroundStation {
	groundStation := &PrecomputedGroundStation{
		BaseNode:     BaseNode{Name: name, Router: router, Computing: computing},
		LinkProtocol: linkProtocol,
		positions:    make(map[time.Time]types.Vector),
	}

	router.Mount(groundStation)
	computing.Mount(groundStation)
	linkProtocol.Mount(groundStation)
	return groundStation
}

func (s *PrecomputedGroundStation) UpdatePosition(time time.Time) {
	s.Position = s.positions[time]
}

func (s *PrecomputedGroundStation) GetLinkNodeProtocol() types.LinkNodeProtocol {
	return s.LinkProtocol
}

func (s *PrecomputedGroundStation) AddPositionState(time time.Time, position types.Vector) {
	s.positions[time] = position
}
