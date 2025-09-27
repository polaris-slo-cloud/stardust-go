package linktypes

import (
	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.Link = (*PrecomputedLink)(nil)

type PrecomputedLink struct {
	Node1 types.Node
	Node2 types.Node
}

func NewPrecomputedLink(node1 types.Node, node2 types.Node) *PrecomputedLink {
	return &PrecomputedLink{
		Node1: node1,
		Node2: node2,
	}
}

func (l *PrecomputedLink) Distance() float64 {
	return l.Node1.DistanceTo(l.Node2)
}

func (l *PrecomputedLink) Latency() float64 {
	return l.Distance() / linkSpeed * 1000
}

func (l *PrecomputedLink) Bandwidth() float64 {
	return 200_000_000_000 // 200 Gbps
}

func (l *PrecomputedLink) IsReachable() bool {
	v := l.Node2.GetPosition().Subtract(l.Node1.GetPosition())
	cross := v.Cross(l.Node1.GetPosition())
	d := cross.Magnitude() / v.Magnitude()
	return d > configs.EarthRadius+10_000 // 10 km buffer
}

func (l *PrecomputedLink) GetOther(self types.Node) types.Node {
	if self == l.Node1 {
		return l.Node2
	}
	if self == l.Node2 {
		return l.Node1
	}
	return nil
}

func (l *PrecomputedLink) Nodes() (types.Node, types.Node) {
	return l.Node1, l.Node2
}
