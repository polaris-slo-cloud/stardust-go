package node

import (
	"github.com/keniack/stardustGo/pkg/types"
	"math"
)

// Node defines a common interface/base for satellites and ground stations
// Since Go has no abstract classes, we use interface + embedding

type Node struct {
	Name      string
	Router    types.IRouter
	Computing types.IComputing
	Position  types.Vector
}

func (n *Node) DistanceTo(other *Node) float64 {
	dx := other.Position.X - n.Position.X
	dy := other.Position.Y - n.Position.Y
	dz := other.Position.Z - n.Position.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (n *Node) GetComputing() types.IComputing {
	return n.Computing
}
