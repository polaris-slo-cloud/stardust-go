package types

import "time"

// Node represents any node in the simulation (satellite or ground).
type Node interface {

	// GetName returns the name of the node
	GetName() string

	// GetRouter returns the router instance associated with the node
	GetRouter() Router

	// GetComputing returns the computing instance associated with the node
	GetComputing() Computing

	// GetPosition returns the current position of the node in ECEF coordinates
	GetPosition() Vector

	// DistanceTo computes the distance to another node in meters
	DistanceTo(other Node) float64

	// UpdatePosition updates the node's position based on the simulation time
	UpdatePosition(simTime time.Time)

	// GetLinkNodeProtocol returns the link protocol instance associated with the node
	GetLinkNodeProtocol() LinkNodeProtocol
}
