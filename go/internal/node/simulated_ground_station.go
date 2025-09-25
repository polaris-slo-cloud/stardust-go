package node

import (
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.GroundStation = (*SimualtedGroundStation)(nil)

type SimualtedGroundStation struct {
	BaseNode
	LinkProtocol types.LinkNodeProtocol
}

func NewSimulatedGroundStation(name string, router types.Router, computing types.Computing, linkProtocol types.LinkNodeProtocol) *SimualtedGroundStation {
	return &SimualtedGroundStation{
		BaseNode:     BaseNode{Name: name, Router: router, Computing: computing},
		LinkProtocol: linkProtocol,
	}
}

func (s *SimualtedGroundStation) UpdatePosition(time time.Time) {
	// no-op for simulated ground stations
}

func (s *SimualtedGroundStation) GetLinkNodeProtocol() types.LinkNodeProtocol {
	return s.LinkProtocol
}
