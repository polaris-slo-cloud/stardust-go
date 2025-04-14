package links

import (
	"errors"
	"sync"

	"stardustGo/internal/links/linktypes"
	"stardustGo/pkg/types"
)

// IslFilterProtocol wraps another IInterSatelliteLinkProtocol and filters links
// so that only those involving the mounted node are retained and processed.
type IslFilterProtocol struct {
	inner       types.IInterSatelliteLinkProtocol
	satellite   types.INode
	links       map[*linktypes.IslLink]struct{}
	established map[*linktypes.IslLink]struct{}
	mu          sync.Mutex
}

// NewIslFilterProtocol initializes the filter protocol
func NewIslFilterProtocol(inner types.IInterSatelliteLinkProtocol) *IslFilterProtocol {
	return &IslFilterProtocol{
		inner:       inner,
		links:       make(map[*linktypes.IslLink]struct{}),
		established: make(map[*linktypes.IslLink]struct{}),
	}
}

// Mount binds the protocol to a specific node
func (p *IslFilterProtocol) Mount(s types.INode) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.satellite == nil {
		p.satellite = s
		p.inner.Mount(s)
	}
}

// AddLink includes a link if relevant to the mounted node
func (p *IslFilterProtocol) AddLink(link types.ILink) {
	if isl, ok := link.(*linktypes.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		if isl.Involves(p.satellite) {
			p.links[isl] = struct{}{}
		}
		p.inner.AddLink(isl)
	}
}

// ConnectLink establishes a specific link if relevant
func (p *IslFilterProtocol) ConnectLink(link types.ILink) error {
	if isl, ok := link.(*linktypes.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		if _, ok := p.links[isl]; ok {
			p.established[isl] = struct{}{}
			isl.SetEstablished(true)
		}
		return p.inner.ConnectLink(isl)
	}
	return nil
}

// DisconnectLink removes a link from the established set
func (p *IslFilterProtocol) DisconnectLink(link types.ILink) error {
	if isl, ok := link.(*linktypes.IslLink); ok {
		p.mu.Lock()
		defer p.mu.Unlock()
		delete(p.established, isl)
		isl.SetEstablished(false)
		return p.inner.DisconnectLink(isl)
	}
	return nil
}

// ConnectSatellite connects to all links involving the given satellite
func (p *IslFilterProtocol) ConnectSatellite(s types.INode) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if s.GetName() == p.satellite.GetName() {
		return errors.New("cannot connect to self")
	}
	for l := range p.links {
		if l.Involves(s) {
			_ = p.ConnectLink(l)
		}
	}
	return nil
}

// DisconnectSatellite disconnects all links involving the given satellite
func (p *IslFilterProtocol) DisconnectSatellite(s types.INode) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if s.GetName() == p.satellite.GetName() {
		return errors.New("cannot disconnect self")
	}
	for l := range p.links {
		if l.Involves(s) {
			_ = p.DisconnectLink(l)
		}
	}
	return nil
}

// UpdateLinks applies the inner protocol update and filters results
func (p *IslFilterProtocol) UpdateLinks() ([]types.ILink, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.satellite == nil {
		return nil, errors.New("not mounted")
	}

	all, err := p.inner.UpdateLinks()
	if err != nil {
		return nil, err
	}

	filtered := make([]types.ILink, 0, len(all))
	for _, link := range all {
		if isl, ok := link.(*linktypes.IslLink); ok {
			if isl.Involves(p.satellite) {
				filtered = append(filtered, isl)
				p.established[isl] = struct{}{}
			}
		}
	}

	return filtered, nil
}

// Links returns all relevant links for the mounted node
func (p *IslFilterProtocol) Links() []types.ILink {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]types.ILink, 0, len(p.links))
	for l := range p.links {
		out = append(out, l)
	}
	return out
}

// Established returns only active links involving the node
func (p *IslFilterProtocol) Established() []types.ILink {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]types.ILink, 0, len(p.established))
	for l := range p.established {
		out = append(out, l)
	}
	return out
}
