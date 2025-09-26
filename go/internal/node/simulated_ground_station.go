package node

import (
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.GroundStation = (*SimulatedGroundStation)(nil)

type SimulatedGroundStation struct {
	BaseNode

	LinkProtocol types.LinkNodeProtocol
	positions    map[time.Time]types.Vector
}

func NewSimulatedGroundStation(name string, router types.Router, computing types.Computing, linkProtocol types.LinkNodeProtocol) *SimulatedGroundStation {
	groundStation := &SimulatedGroundStation{
		BaseNode:     BaseNode{Name: name, Router: router, Computing: computing},
		LinkProtocol: linkProtocol,
		positions:    make(map[time.Time]types.Vector),
	}

	router.Mount(groundStation)
	computing.Mount(groundStation)
	linkProtocol.Mount(groundStation)
	return groundStation
}

func (s *SimulatedGroundStation) UpdatePosition(time time.Time) {
	s.Position = s.positions[time]
}

func (s *SimulatedGroundStation) GetLinkNodeProtocol() types.LinkNodeProtocol {
	return s.LinkProtocol
}

func (s *SimulatedGroundStation) AddPositionState(time time.Time, position types.Vector) {
	s.positions[time] = position
}
