package types

// Satellite represents a satellite node
type Satellite interface {
	Node

	// GetISLProtocol returns the ISL protocol
	GetISLProtocol() InterSatelliteLinkProtocol
}
