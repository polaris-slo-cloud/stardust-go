package simulation

import (
	"time"

	"github.com/keniack/stardustGo/internal/node"
	"github.com/keniack/stardustGo/pkg/types"
)

type SimulationController interface {
	InjectSatellites([]types.Node) error
	InjectGroundStations([]types.Node) error
	StartAutorun() <-chan struct{}
	StopAutorun()
	StepBySeconds(seconds float64)
	StepByTime(newTime time.Time)
	GetAllNodes() []types.Node
	GetSatellites() []*node.Satellite
	GetGroundStations() []*node.GroundStation
}
