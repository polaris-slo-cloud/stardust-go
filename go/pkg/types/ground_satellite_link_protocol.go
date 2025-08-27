package types

// IGroundSatelliteLinkProtocol abstracts ground link handling logic
type IGroundSatelliteLinkProtocol interface {
	Link() *Link
	UpdateLink() error
	Mount(station *Node)
}
