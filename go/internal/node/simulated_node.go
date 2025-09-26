package node

import (
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ SimulatedNode = (*SimulatedSatellite)(nil)
var _ SimulatedNode = (*SimulatedGroundStation)(nil)

type SimulatedNode interface {
	types.Node
	AddPositionState(time time.Time, position types.Vector)
}
