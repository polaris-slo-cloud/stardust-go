package links

import (
	"errors"
	"sort"
	"sync"

	"github.com/keniack/stardustGo/internal/links/linktypes"
	"github.com/keniack/stardustGo/pkg/types"
)

// GroundSatelliteNearestProtocol maintains a single active link from the ground station
// to the nearest satellite at any given time.
type GroundSatelliteNearestProtocol struct {
	link          *linktypes.GroundLink // Current active ground link
	satellites    []types.INode         // Available satellites
	groundStation types.INode           // The ground station node
	mu            sync.Mutex
}

// NewGroundSatelliteNearestProtocol creates a new protocol with an initial list of satellites.
func NewGroundSatelliteNearestProtocol(satellites []types.INode) *GroundSatelliteNearestProtocol {
	return &GroundSatelliteNearestProtocol{
		satellites: satellites,
	}
}

// Mount binds this protocol to a ground station.
func (p *GroundSatelliteNearestProtocol) Mount(gs types.INode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.groundStation == nil {
		p.groundStation = gs
	}
}

// AddLink is a no-op for this protocol.
func (p *GroundSatelliteNearestProtocol) AddLink(link types.ILink) {}

// ConnectLink is a no-op for this protocol.
func (p *GroundSatelliteNearestProtocol) ConnectLink(link types.ILink) error {
	return nil
}

// DisconnectLink is a no-op for this protocol.
func (p *GroundSatelliteNearestProtocol) DisconnectLink(link types.ILink) error {
	return nil
}

// ConnectSatellite is not used in this context.
func (p *GroundSatelliteNearestProtocol) ConnectSatellite(s types.INode) error {
	return nil
}

// DisconnectSatellite is not used in this context.
func (p *GroundSatelliteNearestProtocol) DisconnectSatellite(s types.INode) error {
	return nil
}

// UpdateLinks selects the closest satellite and sets up the ground link accordingly.
func (p *GroundSatelliteNearestProtocol) UpdateLinks() ([]types.ILink, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.groundStation == nil {
		return nil, errors.New("protocol not mounted to ground station")
	}
	if len(p.satellites) == 0 {
		return nil, errors.New("no satellites available")
	}

	sort.Slice(p.satellites, func(i, j int) bool {
		return p.groundStation.DistanceTo(p.satellites[i]) < p.groundStation.DistanceTo(p.satellites[j])
	})

	nearest := p.satellites[0]
	if nearest == nil || (p.link != nil && p.link.Satellite.GetName() == nearest.GetName()) {
		return []types.ILink{p.link}, nil // Already linked to the nearest
	}

	old := p.link
	p.link = linktypes.NewGroundLink(p.groundStation, nearest)

	// Add new link to satellite if it supports ground links
	if s, ok := nearest.(interface{ AddGroundLink(link types.ILink) }); ok {
		s.AddGroundLink(p.link)
	}

	// Remove old link from previous satellite if supported
	if old != nil {
		if oldSat, ok := old.Satellite.(interface{ RemoveGroundLink(station types.INode) }); ok {
			oldSat.RemoveGroundLink(p.groundStation)
		}
	}

	return []types.ILink{p.link}, nil
}

// Links returns the current active link if any.
func (p *GroundSatelliteNearestProtocol) Links() []types.ILink {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.link != nil {
		return []types.ILink{p.link}
	}
	return nil
}

// Established returns the current active link if any.
func (p *GroundSatelliteNearestProtocol) Established() []types.ILink {
	return p.Links()
}

// Link returns the currently active GroundLink.
func (p *GroundSatelliteNearestProtocol) Link() *linktypes.GroundLink {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.link
}

// AddSatellite adds a satellite to the trackable list.
func (p *GroundSatelliteNearestProtocol) AddSatellite(sat types.INode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.satellites = append(p.satellites, sat)
}

// RemoveSatellite removes a satellite from the list and resets the link if needed.
func (p *GroundSatelliteNearestProtocol) RemoveSatellite(sat types.INode) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Filter out the satellite
	filtered := make([]types.INode, 0, len(p.satellites))
	for _, s := range p.satellites {
		if s.GetName() != sat.GetName() {
			filtered = append(filtered, s)
		}
	}
	p.satellites = filtered

	// Remove the link if it's pointing to the removed satellite
	if p.link != nil && p.link.Satellite.GetName() == sat.GetName() {
		if removable, ok := sat.(interface{ RemoveGroundLink(types.INode) }); ok {
			removable.RemoveGroundLink(p.groundStation)
		}
		p.link = nil
	}
}
