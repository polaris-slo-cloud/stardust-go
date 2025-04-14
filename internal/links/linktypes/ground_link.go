package linktypes

import (
	"stardustGo/configs"
	"stardustGo/pkg/types"
)

const groundSpeedOfLight = configs.SpeedOfLight * 0.98 // 98% of light speed

type GroundLink struct {
	GroundStation types.INode
	Satellite     types.INode
}

// NewGroundLink constructs a link between a ground station and a satellite.
func NewGroundLink(gs types.INode, sat types.INode) *GroundLink {
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
	return gl.Distance() / groundSpeedOfLight * 1000
}

// Bandwidth returns the link bandwidth in bits per second.
func (gl *GroundLink) Bandwidth() float64 {
	return 500_000_000 // 500 Mbps
}

// Established always returns true for ground links.
func (gl *GroundLink) Established() bool {
	return true
}

func (gl *GroundLink) GetOther(self types.INode) types.INode {
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
