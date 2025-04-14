package types

import (
	"time"
)

type NodeWithISL interface {
	INode
	InterSatelliteLinkProtocol() IInterSatelliteLinkProtocol
}

// INode represents any node in the simulation (satellite or ground).
type INode interface {
	GetComputing() IComputing
	GetName() string
	PositionVector() Vector
	DistanceTo(other INode) float64
	GetLinks() []ILink
	UpdatePosition(simTime time.Time)
}

// ILink represents a generic network link.
type ILink interface {
	// Distance returns the link distance in meters.
	Distance() float64

	// Latency returns the link latency in milliseconds.
	Latency() float64

	// Bandwidth returns the bandwidth in bits per second.
	Bandwidth() float64

	// Established returns whether the link is currently active.
	Established() bool

	// GetOther returns the opposite node from the provided one.
	GetOther(self INode) INode

	// IsReachable returns true if the link is physically/line-of-sight reachable.
	IsReachable() bool
}

type IInterSatelliteLinkProtocol interface {
	Links() []ILink
	Established() []ILink
	UpdateLinks() ([]ILink, error)
	ConnectSatellite(s INode) error
	ConnectLink(link ILink) error
	DisconnectSatellite(s INode) error
	DisconnectLink(link ILink) error
	Mount(s INode)
	AddLink(link ILink)
}

// IGroundSatelliteLinkProtocol abstracts ground link handling logic
type IGroundSatelliteLinkProtocol interface {
	Link() *ILink
	UpdateLink() error
	Mount(station *INode)
}
