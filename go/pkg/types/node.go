package types

import "time"

// INode represents any node in the simulation (satellite or ground).
type Node interface {
	GetComputing() Computing
	GetName() string
	PositionVector() Vector
	DistanceTo(other Node) float64
	UpdatePosition(simTime time.Time)
	GetLinks() []Link
	GetEstablishedLinks() []Link
	GetLinkNodeProtocol() LinkNodeProtocol
}
