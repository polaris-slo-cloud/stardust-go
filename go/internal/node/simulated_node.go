package node

import (
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

var _ PrecomputedNode = (*PrecomputedSatellite)(nil)
var _ PrecomputedNode = (*PrecomputedGroundStation)(nil)

type PrecomputedNode interface {
	types.Node
	AddPositionState(time time.Time, position types.Vector)
}
