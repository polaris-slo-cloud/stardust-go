package links

import (
	"sync"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/internal/links/linktypes"
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

// ConnectSatellite delegates to the wrapped protocol.
func (p *IslAddLoopProtocol) ConnectSatellite(s types.Node) error {
	return p.inner.ConnectSatellite(s)
}

// DisconnectSatellite delegates to the wrapped protocol.
func (p *IslAddLoopProtocol) DisconnectSatellite(s types.Node) error {
	return p.inner.DisconnectSatellite(s)
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

	var established []*linktypes.IslLink
	for _, l := range innerEstablished {
		if isl, ok := l.(*linktypes.IslLink); ok {
			established = append(established, isl)
		}
	}

	// Only add an extra link if we're under the target neighbor count
	if len(established) > 0 && len(established) < p.config.Neighbours-1 {
		type candidate struct {
			dist float64
			link *linktypes.IslLink
		}
		var best *candidate

		for _, l := range p.inner.Links() {
			isl, ok := l.(*linktypes.IslLink)
			if !ok || contains(established, isl) || isl.Distance() > configs.MaxISLDistance {
				continue
			}

			n1, ok1 := isl.Node1.(types.NodeWithISL)
			n2, ok2 := isl.Node2.(types.NodeWithISL)
			if !ok1 || !ok2 {
				continue
			}

			if !shouldLoop(n1.InterSatelliteLinkProtocol().Established(), p.config.Neighbours) ||
				!shouldLoop(n2.InterSatelliteLinkProtocol().Established(), p.config.Neighbours) {
				continue
			}

			if best == nil || isl.Distance() < best.dist {
				best = &candidate{dist: isl.Distance(), link: isl}
			}
		}

		if best != nil {
			n1 := best.link.Node1.(types.NodeWithISL)
			n2 := best.link.Node2.(types.NodeWithISL)
			_ = n1.InterSatelliteLinkProtocol().ConnectLink(best.link)
			_ = n2.InterSatelliteLinkProtocol().ConnectLink(best.link)
			best.link.SetEstablished(true)
			established = append(established, best.link)
		}
	}

	// Return as []types.ILink
	out := make([]types.Link, len(established))
	for i, l := range established {
		out[i] = l
	}
	return out, nil
}

// contains checks if a link is already in the list.
func contains(list []*linktypes.IslLink, link *linktypes.IslLink) bool {
	for _, l := range list {
		if l == link {
			return true
		}
	}
	return false
}

// shouldLoop returns true if the node has fewer than the maximum allowed neighbors.
func shouldLoop(established []types.Link, max int) bool {
	return len(established) < max
}
