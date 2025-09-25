package linktypes

import (
	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.Link = (*SimulatedLink)(nil)

type SimulatedLink struct {
	Node1 types.Node
	Node2 types.Node
}

func NewSimulatedLink(node1 types.Node, node2 types.Node) *SimulatedLink {
	return &SimulatedLink{
		Node1: node1,
		Node2: node2,
	}
}

func (l *SimulatedLink) Distance() float64 {
	return l.Node1.DistanceTo(l.Node2)
}

func (l *SimulatedLink) Latency() float64 {
	return l.Distance() / speedOfLight * 1000
}

func (l *SimulatedLink) Bandwidth() float64 {
	return 200_000_000_000 // 200 Gbps
}

func (l *SimulatedLink) IsReachable() bool {
	v := l.Node2.PositionVector().Subtract(l.Node1.PositionVector())
	cross := v.Cross(l.Node1.PositionVector())
	d := cross.Magnitude() / v.Magnitude()
	return d > configs.EarthRadius+10_000 // 10 km buffer
}

func (l *SimulatedLink) GetOther(self types.Node) types.Node {
	if self.GetName() == l.Node1.GetName() {
		return l.Node2
	}
	if self.GetName() == l.Node2.GetName() {
		return l.Node1
	}
	return nil
}

func (l *SimulatedLink) Nodes() (types.Node, types.Node) {
	return l.Node1, l.Node2
}
