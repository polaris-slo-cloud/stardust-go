package types

import (
	"math"
)

// Node defines a common interface/base for satellites and ground stations
// Since Go has no abstract classes, we use interface + embedding

type BaseNode struct {
	Name      string
	Router    Router
	Computing Computing
	Position  Vector
}

func (n *BaseNode) DistanceTo(other *BaseNode) float64 {
	dx := other.Position.X - n.Position.X
	dy := other.Position.Y - n.Position.Y
	dz := other.Position.Z - n.Position.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (n *BaseNode) GetComputing() Computing {
	return n.Computing
}
