package links

import (
	"errors"
	"sort"
	"sync"

	"github.com/keniack/stardustGo/internal/links/linktypes"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.GroundSatelliteLinkProtocol = (*GroundSatelliteNearestProtocol)(nil)

// GroundSatelliteNearestProtocol maintains a single active link from the ground station
// to the nearest satellite at any given time.
type GroundSatelliteNearestProtocol struct {
	link          *linktypes.GroundLink // Current active ground link
	satellites    []types.Node          // Available satellites
	groundStation types.Node            // The ground station node
	mu            sync.Mutex
}

// NewGroundSatelliteNearestProtocol creates a new protocol with an initial list of satellites.
func NewGroundSatelliteNearestProtocol(satellites []types.Node) *GroundSatelliteNearestProtocol {
	return &GroundSatelliteNearestProtocol{
		satellites: satellites,
	}
}

// Mount binds this protocol to a ground station.
func (p *GroundSatelliteNearestProtocol) Mount(gs types.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.groundStation == nil {
		p.groundStation = gs
	}
}

// AddLink is a no-op for this protocol.
func (p *GroundSatelliteNearestProtocol) AddLink(link types.Link) {}

// ConnectLink is a no-op for this protocol.
func (p *GroundSatelliteNearestProtocol) ConnectLink(link types.Link) error {
	return nil
}

// DisconnectLink is a no-op for this protocol.
func (p *GroundSatelliteNearestProtocol) DisconnectLink(link types.Link) error {
	return nil
}

// ConnectSatellite is not used in this context.
func (p *GroundSatelliteNearestProtocol) ConnectSatellite(s types.Node) error {
	return nil
}

// DisconnectSatellite is not used in this context.
func (p *GroundSatelliteNearestProtocol) DisconnectSatellite(s types.Node) error {
	return nil
}

// UpdateLink selects the closest satellite and sets up the ground link accordingly.
func (p *GroundSatelliteNearestProtocol) UpdateLink() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.groundStation == nil {
		return errors.New("protocol not mounted to ground station")
	}
	if len(p.satellites) == 0 {
		return errors.New("no satellites available")
	}

	sort.Slice(p.satellites, func(i, j int) bool {
		return p.groundStation.DistanceTo(p.satellites[i]) < p.groundStation.DistanceTo(p.satellites[j])
	})

	nearest := p.satellites[0]
	if nearest == nil || (p.link != nil && p.link.Satellite.GetName() == nearest.GetName()) {
		return nil // Already linked to the nearest
	}

	old := p.link
	p.link = linktypes.NewGroundLink(p.groundStation, nearest)

	// Add new link to satellite if it supports ground links
	if s, ok := nearest.(interface{ AddGroundLink(link types.Link) }); ok {
		s.AddGroundLink(p.link)
	}

	// Remove old link from previous satellite if supported
	if old != nil {
		if oldSat, ok := old.Satellite.(interface{ RemoveGroundLink(station types.Node) }); ok {
			oldSat.RemoveGroundLink(p.groundStation)
		}
	}

	return nil
}

// Links returns the current active link if any.
func (p *GroundSatelliteNearestProtocol) Links() []types.Link {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.link != nil {
		return []types.Link{p.link}
	}
	return nil
}

// Established returns the current active link if any.
func (p *GroundSatelliteNearestProtocol) Established() []types.Link {
	return p.Links()
}

// Link returns the currently active GroundLink.
func (p *GroundSatelliteNearestProtocol) Link() *types.Link {
	p.mu.Lock()
	defer p.mu.Unlock()
	var link types.Link = p.link
	return &link
}

// AddSatellite adds a satellite to the trackable list.
func (p *GroundSatelliteNearestProtocol) AddSatellite(sat types.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.satellites = append(p.satellites, sat)
}

// RemoveSatellite removes a satellite from the list and resets the link if needed.
func (p *GroundSatelliteNearestProtocol) RemoveSatellite(sat types.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Filter out the satellite
	filtered := make([]types.Node, 0, len(p.satellites))
	for _, s := range p.satellites {
		if s.GetName() != sat.GetName() {
			filtered = append(filtered, s)
		}
	}
	p.satellites = filtered

	// Remove the link if it's pointing to the removed satellite
	if p.link != nil && p.link.Satellite.GetName() == sat.GetName() {
		if removable, ok := sat.(interface{ RemoveGroundLink(types.Node) }); ok {
			removable.RemoveGroundLink(p.groundStation)
		}
		p.link = nil
	}
}
