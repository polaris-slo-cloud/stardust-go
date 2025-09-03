package types

// GroundSatelliteLinkProtocol abstracts ground link handling logic
type GroundSatelliteLinkProtocol interface {
	Link() Link
	UpdateLink() error
	Mount(station Node)
}
