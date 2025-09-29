package linktypes

import (
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.Link = (*GroundLink)(nil)

type GroundLink struct {
	GroundStation types.Node
	Satellite     types.Node
}

// NewGroundLink constructs a link between a ground station and a satellite.
func NewGroundLink(gs types.Node, sat types.Node) *GroundLink {
	return &GroundLink{
		GroundStation: gs,
		Satellite:     sat,
	}
}

// Distance returns the distance in meters between the ground station and satellite.
func (gl *GroundLink) Distance() float64 {
	return gl.GroundStation.DistanceTo(gl.Satellite)
}

// Latency returns the one-way latency in milliseconds.
func (gl *GroundLink) Latency() float64 {
	return gl.Distance() / linkSpeed * 1000
}

// Bandwidth returns the link bandwidth in bits per second.
func (gl *GroundLink) Bandwidth() float64 {
	return 500_000_000 // 500 Mbps
}

func (gl *GroundLink) GetOther(self types.Node) types.Node {
	if self.GetName() == gl.Satellite.GetName() {
		return gl.GroundStation
	}
	if self.GetName() == gl.GroundStation.GetName() {
		return gl.Satellite
	}
	// Return nil or panic, depending on how you want to fail
	return nil
}

// IsReachable returns true – placeholder for future visibility checks.
func (gl *GroundLink) IsReachable() bool {
	return true
}

func (gl *GroundLink) Nodes() (types.Node, types.Node) {
	return gl.GroundStation, gl.Satellite
}
