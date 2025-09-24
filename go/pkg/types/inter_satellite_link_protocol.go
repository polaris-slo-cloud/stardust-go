package types

// InterSatelliteLinkProtocol defines the interface for managing inter-satellite links (ISL) between nodes.
type InterSatelliteLinkProtocol interface {
	LinkNodeProtocol

	ConnectSatellite(s Node) error
	DisconnectSatellite(s Node) error
	AddLink(link Link)
}
