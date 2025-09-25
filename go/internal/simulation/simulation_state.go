package simulation

import (
	"time"

	"github.com/keniack/stardustGo/pkg/types"
)

type SimulationLink struct {
	NodeName1 string
	NodeName2 string
}

type SimulationMetadata struct {
	Satellites []RawSatellite
	Grounds    []RawGroundStation
	Links      []SimulationLink
	States     []SimulationState
}

type SimulationState struct {
	Time       time.Time
	NodeStates []NodeState
}

type NodeState struct {
	Name        string
	Position    types.Vector
	Established []int
}

type RawSatellite struct {
	Name          string
	ComputingType types.ComputingType
}

type RawGroundStation struct {
	Name          string
	ComputingType types.ComputingType
}

func NewSimulationMetadata() SimulationMetadata {
	return SimulationMetadata{
		Links:  []SimulationLink{},
		States: []SimulationState{},
	}
}

func NewSimulationState(time time.Time, nodes []NodeState) SimulationState {
	return SimulationState{
		Time:       time,
		NodeStates: nodes,
	}
}

func NewNodeState(name string, position types.Vector, established []int) NodeState {
	return NodeState{
		Name:        name,
		Position:    position,
		Established: established,
	}
}
