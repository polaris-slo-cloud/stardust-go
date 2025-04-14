package linktypes

import (
	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
)

const speedOfLight = configs.SpeedOfLight * 0.99 // 99% of light speed

// IslLink represents an inter-satellite laser link.
type IslLink struct {
	Node1         types.INode
	Node2         types.INode
	isEstablished bool
}

// NewIslLink creates a new ISL between two nodes.
func NewIslLink(n1, n2 types.INode) *IslLink {
	return &IslLink{
		Node1: n1,
		Node2: n2,
	}
}

// Distance returns the link distance in meters.
func (l *IslLink) Distance() float64 {
	pos1 := l.Node1.PositionVector()
	pos2 := l.Node2.PositionVector()
	return pos1.Sub(pos2).Magnitude()
}

// Latency returns the communication latency in milliseconds.
func (l *IslLink) Latency() float64 {
	return l.Distance() / speedOfLight * 1000
}

// Bandwidth returns the bandwidth in bits per second.
func (l *IslLink) Bandwidth() float64 {
	return 200_000_000_000 // 200 Gbps
}

// IsReachable checks if line-of-sight is available.
func (l *IslLink) IsReachable() bool {
	v := l.Node2.PositionVector().Sub(l.Node1.PositionVector())
	cross := v.Cross(l.Node1.PositionVector())
	d := cross.Magnitude() / v.Magnitude()
	return d > configs.EarthRadius+10_000 // 10 km buffer
}

// GetOther returns the opposite node of the link.
func (l *IslLink) GetOther(self types.INode) types.INode {
	if self.GetName() == l.Node1.GetName() {
		return l.Node2
	}
	if self.GetName() == l.Node2.GetName() {
		return l.Node1
	}
	return nil
}

func (l *IslLink) Involves(node types.INode) bool {
	return l.Node1.GetName() == node.GetName() || l.Node2.GetName() == node.GetName()
}

func (l *IslLink) Established() bool {
	return l.isEstablished
}

func (l *IslLink) SetEstablished(val bool) {
	l.isEstablished = val
}
