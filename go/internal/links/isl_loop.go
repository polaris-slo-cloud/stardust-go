package links

import (
	"slices"
	"sync"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.InterSatelliteLinkProtocol = (*IslAddLoopProtocol)(nil)

// IslAddLoopProtocol wraps another IInterSatelliteLinkProtocol
// and adds a single additional loop link if too few are established.
type IslAddLoopProtocol struct {
	inner  types.InterSatelliteLinkProtocol
	config configs.InterSatelliteLinkConfig
	mu     sync.Mutex
}

// NewIslAddLoopProtocol creates a loop-adding decorator.
func NewIslAddLoopProtocol(inner types.InterSatelliteLinkProtocol, cfg configs.InterSatelliteLinkConfig) *IslAddLoopProtocol {
	return &IslAddLoopProtocol{inner: inner, config: cfg}
}

// Mount delegates mounting to the wrapped protocol.
func (p *IslAddLoopProtocol) Mount(s types.Node) {
	p.inner.Mount(s)
}

// AddLink delegates link registration to the wrapped protocol.
func (p *IslAddLoopProtocol) AddLink(link types.Link) {
	p.inner.AddLink(link)
}

// ConnectLink delegates connection to the wrapped protocol.
func (p *IslAddLoopProtocol) ConnectLink(link types.Link) error {
	return p.inner.ConnectLink(link)
}

// DisconnectLink delegates disconnection to the wrapped protocol.
func (p *IslAddLoopProtocol) DisconnectLink(link types.Link) error {
	return p.inner.DisconnectLink(link)
}

// Links returns all *IslLink links from the wrapped protocol.
func (p *IslAddLoopProtocol) Links() []types.Link {
	return p.inner.Links()
}

// Established returns all active *IslLink connections.
func (p *IslAddLoopProtocol) Established() []types.Link {
	return p.inner.Established()
}

// UpdateLinks adds one additional loop link if there are too few established connections.
func (p *IslAddLoopProtocol) UpdateLinks() ([]types.Link, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	innerEstablished, err := p.inner.UpdateLinks()
	if err != nil {
		return nil, err
	}

	established := make([]types.Link, len(innerEstablished))
	copy(established, innerEstablished)

	// Only add an extra link if we're under the target neighbor count
	if len(established) > 0 && len(established) < p.config.Neighbours-1 {
		type candidate struct {
			dist float64
			link types.Link
		}
		var best *candidate

		for _, l := range p.inner.Links() {
			if slices.Contains(established, l) || l.Distance() > configs.MaxISLDistance {
				continue
			}

			n1, n2 := l.Nodes()
			if !shouldLoop(n1.GetLinkNodeProtocol().Established(), p.config.Neighbours) ||
				!shouldLoop(n2.GetLinkNodeProtocol().Established(), p.config.Neighbours) {
				continue
			}

			if best == nil || l.Distance() < best.dist {
				best = &candidate{dist: l.Distance(), link: l}
			}
		}

		if best != nil {
			n1, n2 := best.link.Nodes()
			_ = n1.GetLinkNodeProtocol().ConnectLink(best.link)
			_ = n2.GetLinkNodeProtocol().ConnectLink(best.link)
			established = append(established, best.link)
		}
	}

	return established, nil
}

// shouldLoop returns true if the node has fewer than the maximum allowed neighbors.
func shouldLoop(established []types.Link, max int) bool {
	return len(established) < max
}
