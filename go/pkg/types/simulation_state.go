package types

import (
	"time"
)

type SimulationLink struct {
	NodeName1 string
	NodeName2 string
}

type SimulationMetadata struct {
	StatePlugins []string
	Satellites   []RawSatellite
	Grounds      []RawGroundStation
	Links        []SimulationLink
	States       []SimulationState
}

type SimulationState struct {
	Time       time.Time
	NodeStates []NodeState
}

type NodeState struct {
	Name        string
	Position    Vector
	Established []int
}

type RawSatellite struct {
	Name          string
	ComputingType ComputingType
}

type RawGroundStation struct {
	Name          string
	ComputingType ComputingType
}

func NewSimulationMetadata() SimulationMetadata {
	return SimulationMetadata{
		StatePlugins: []string{},
		Links:        []SimulationLink{},
		States:       []SimulationState{},
	}
}

func NewSimulationState(time time.Time, nodes []NodeState) SimulationState {
	return SimulationState{
		Time:       time,
		NodeStates: nodes,
	}
}

func NewNodeState(name string, position Vector, established []int) NodeState {
	return NodeState{
		Name:        name,
		Position:    position,
		Established: established,
	}
}
