package links

import "github.com/keniack/stardustGo/pkg/types"

var _ types.InterSatelliteLinkProtocol = (*SimulatedLinkProtocol)(nil)

type SimulatedLinkProtocol struct {
}

func NewSimulatedLinkProtocol() *SimulatedLinkProtocol {
	return &SimulatedLinkProtocol{}
}

// Mount attaches the protocol to a node.
func (p *SimulatedLinkProtocol) Mount(s types.Node) {
	// Base implementation: do nothing.
}

// ConnectLink establishes a connection for the given link.
func (p *SimulatedLinkProtocol) ConnectLink(link types.Link) error {
	// Base implementation: return nil (success).
	return nil
}

// DisconnectLink disconnects the given link.
func (p *SimulatedLinkProtocol) DisconnectLink(link types.Link) error {
	// Base implementation: return nil (success).
	return nil
}

// UpdateLinks updates the list of active links.
func (p *SimulatedLinkProtocol) UpdateLinks() ([]types.Link, error) {
	// Base implementation: return empty slice and nil error.
	return []types.Link{}, nil
}

// Established returns the list of established links.
func (p *SimulatedLinkProtocol) Established() []types.Link {
	// Base implementation: return empty slice.
	return []types.Link{}
}

// Links returns all links managed by this protocol.
func (p *SimulatedLinkProtocol) Links() []types.Link {
	// Base implementation: return empty slice.
	return []types.Link{}
}

// ConnectSatellite establishes a connection to a satellite node.
func (p *SimulatedLinkProtocol) ConnectSatellite(s types.Node) error {
	// Base implementation: return nil (success).
	return nil
}

// DisconnectSatellite disconnects from a satellite node.
func (p *SimulatedLinkProtocol) DisconnectSatellite(s types.Node) error {
	// Base implementation: return nil (success).
	return nil
}

// AddLink adds a new link to the protocol.
func (p *SimulatedLinkProtocol) AddLink(link types.Link) {
	// Base implementation: do nothing.
}
