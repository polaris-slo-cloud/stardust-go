package links

import (
	"errors"
	"sort"
	"sync"

	configmod "github.com/keniack/stardustGo/configs"
	linkmod "github.com/keniack/stardustGo/internal/links/linktypes"
	"github.com/keniack/stardustGo/pkg/types"
)

// IslNearestProtocol connects a node to its N nearest neighbors using ISLs.
type IslNearestProtocol struct {
	config    configmod.InterSatelliteLinkConfig
	satellite types.Node

	mu       sync.Mutex
	links    []*linkmod.IslLink        // All potential links
	outgoing []*linkmod.IslLink        // Active outgoing links
	incoming map[*linkmod.IslLink]bool // Remote links initiated by others
}

// NewIslNearestProtocol initializes the nearest-neighbor protocol.
func NewIslNearestProtocol(cfg configmod.InterSatelliteLinkConfig) *IslNearestProtocol {
	return &IslNearestProtocol{
		config:   cfg,
		incoming: make(map[*linkmod.IslLink]bool),
	}
}

// Mount binds the protocol to a satellite node.
func (p *IslNearestProtocol) Mount(s types.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.satellite = s
}

// AddLink registers a new potential inter-satellite link.
func (p *IslNearestProtocol) AddLink(link types.Link) {
	if isl, ok := link.(*linkmod.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.links = append(p.links, isl)
	}
}

// ConnectLink marks an incoming connection from a peer.
func (p *IslNearestProtocol) ConnectLink(link types.Link) error {
	if isl, ok := link.(*linkmod.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.incoming[isl] = true
	}
	return nil
}

// DisconnectLink removes the incoming status if it's not also an outgoing link.
func (p *IslNearestProtocol) DisconnectLink(link types.Link) error {
	if isl, ok := link.(*linkmod.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		delete(p.incoming, isl)
		if !p.isInOutgoing(isl) {
			isl.SetEstablished(false)
		}
	}
	return nil
}

// ConnectSatellite is a helper to find and connect to a specific satellite.
func (p *IslNearestProtocol) ConnectSatellite(s types.Node) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, l := range p.links {
		if l.Node1 == s || l.Node2 == s {
			return p.ConnectLink(l)
		}
	}
	return errors.New("no link to target satellite found")
}

// DisconnectSatellite is a helper to find and disconnect from a specific satellite.
func (p *IslNearestProtocol) DisconnectSatellite(s types.Node) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, l := range p.links {
		if l.Node1 == s || l.Node2 == s {
			return p.DisconnectLink(l)
		}
	}
	return nil
}

// UpdateLinks reconnects to the closest N reachable satellites.
func (p *IslNearestProtocol) UpdateLinks() ([]types.Link, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.satellite == nil {
		return nil, errors.New("protocol not mounted")
	}

	// Track previous outgoing links
	prevOut := make(map[*linkmod.IslLink]struct{})
	for _, l := range p.outgoing {
		prevOut[l] = struct{}{}
	}

	// Collect reachable links
	valid := []*linkmod.IslLink{}
	for _, l := range p.links {
		if l.IsReachable() {
			valid = append(valid, l)
		}
	}

	// Sort by distance
	sort.Slice(valid, func(i, j int) bool {
		return valid[i].Distance() < valid[j].Distance()
	})

	// Select top-N
	selected := valid
	if len(selected) > p.config.Neighbours {
		selected = selected[:p.config.Neighbours]
	}
	p.outgoing = selected

	// Establish new links
	for _, l := range selected {
		if _, seen := prevOut[l]; !seen {
			if other := l.GetOther(p.satellite); other != nil {
				if islNode, ok := other.(types.NodeWithISL); ok {
					_ = islNode.InterSatelliteLinkProtocol().ConnectLink(l)
				}
			}
			l.SetEstablished(true)
		} else {
			delete(prevOut, l)
		}
	}

	// Disconnect dropped links
	for l := range prevOut {
		if other := l.GetOther(p.satellite); other != nil {
			if islNode, ok := other.(types.NodeWithISL); ok {
				_ = islNode.InterSatelliteLinkProtocol().DisconnectLink(l)
			}
		}
		l.SetEstablished(false)
	}

	// Return current active links
	result := make([]types.Link, len(p.outgoing))
	for i, l := range p.outgoing {
		result[i] = l
	}
	return result, nil
}

// Links returns all known links.
func (p *IslNearestProtocol) Links() []types.Link {
	p.mu.Lock()
	defer p.mu.Unlock()
	res := make([]types.Link, len(p.links))
	for i, l := range p.links {
		res[i] = l
	}
	return res
}

// Established returns all active links (incoming or outgoing).
func (p *IslNearestProtocol) Established() []types.Link {
	p.mu.Lock()
	defer p.mu.Unlock()
	seen := make(map[*linkmod.IslLink]bool)
	for _, l := range p.outgoing {
		seen[l] = true
	}
	for l := range p.incoming {
		seen[l] = true
	}
	out := make([]types.Link, 0, len(seen))
	for l := range seen {
		out = append(out, l)
	}
	return out
}

// isInOutgoing checks if a link is in the current outgoing set.
func (p *IslNearestProtocol) isInOutgoing(link *linkmod.IslLink) bool {
	for _, l := range p.outgoing {
		if l == link {
			return true
		}
	}
	return false
}
