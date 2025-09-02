package types

// InterSatelliteLinkProtocol defines the interface for managing inter-satellite links (ISL) between nodes.
type InterSatelliteLinkProtocol interface {
	Links() []Link
	Established() []Link
	UpdateLinks() ([]Link, error)
	ConnectSatellite(s Node) error
	ConnectLink(link Link) error
	DisconnectSatellite(s Node) error
	DisconnectLink(link Link) error
	Mount(s Node)
	AddLink(link Link)
}
