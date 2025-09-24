package node

import (
	"math"

	"github.com/keniack/stardustGo/pkg/types"
)

// Node defines a common interface/base for satellites and ground stations
// Since Go has no abstract classes, we use interface + embedding

type BaseNode struct {
	Name      string
	Router    types.Router
	Computing types.Computing
	Position  types.Vector
}

func (n *BaseNode) DistanceTo(other *BaseNode) float64 {
	dx := other.Position.X - n.Position.X
	dy := other.Position.Y - n.Position.Y
	dz := other.Position.Z - n.Position.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (n *BaseNode) GetComputing() types.Computing {
	return n.Computing
}
