package types

type Satellite interface {
	Node

	GetISLProtocol() InterSatelliteLinkProtocol
}
