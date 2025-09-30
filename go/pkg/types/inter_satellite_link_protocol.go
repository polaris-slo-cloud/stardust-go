package types

// InterSatelliteLinkProtocol defines the interface for managing inter-satellite links (ISL) between nodes.
type InterSatelliteLinkProtocol interface {
	LinkNodeProtocol

	// AddLink adds a new link to the protocol's management
	AddLink(link Link)
}
