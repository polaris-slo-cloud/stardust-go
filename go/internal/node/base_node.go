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

func (n *BaseNode) GetName() string {
	return n.Name
}

func (n *BaseNode) GetRouter() types.Router {
	return n.Router
}

func (n *BaseNode) GetComputing() types.Computing {
	return n.Computing
}

func (n *BaseNode) GetPosition() types.Vector {
	return n.Position
}

func (n *BaseNode) DistanceTo(other types.Node) float64 {
	dx := other.GetPosition().X - n.Position.X
	dy := other.GetPosition().Y - n.Position.Y
	dz := other.GetPosition().Z - n.Position.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
