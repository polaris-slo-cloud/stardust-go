package links

import (
	"errors"
	"sync"

	"github.com/keniack/stardustGo/configs"
	"github.com/keniack/stardustGo/pkg/helper"
	"github.com/keniack/stardustGo/pkg/types"
)

var _ types.InterSatelliteLinkProtocol = (*IslAddSmartLoopProtocol)(nil)

// IslAddSmartLoopProtocol wraps another IInterSatelliteLinkProtocol and augments it
// by adding extra links to form smart loops (e.g., when a satellite has only 1 connection).
type IslAddSmartLoopProtocol struct {
	inner       types.InterSatelliteLinkProtocol
	config      configs.InterSatelliteLinkConfig
	satellite   types.Node
	position    types.Vector
	mu          sync.Mutex
	resultCache []types.Link
	resetEvent  *helper.ManualResetEvent // Signals when ready for reuse
}

// NewIslAddSmartLoopProtocol creates a new smart loop-enhancing protocol
func NewIslAddSmartLoopProtocol(inner types.InterSatelliteLinkProtocol, cfg configs.InterSatelliteLinkConfig) *IslAddSmartLoopProtocol {
	return &IslAddSmartLoopProtocol{
		inner:      inner,
		config:     cfg,
		resetEvent: helper.NewManualResetEvent(true),
	}
}

// Mount binds the protocol to a satellite node
func (p *IslAddSmartLoopProtocol) Mount(s types.Node) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.satellite == nil {
		p.satellite = s
		p.inner.Mount(s)
	}
}

// AddLink registers a link to the underlying protocol
func (p *IslAddSmartLoopProtocol) AddLink(link types.Link) {
	p.inner.AddLink(link)
}

// ConnectLink forwards link connection to the inner protocol
func (p *IslAddSmartLoopProtocol) ConnectLink(link types.Link) error {
	return p.inner.ConnectLink(link)
}

// DisconnectLink forwards link disconnection to the inner protocol
func (p *IslAddSmartLoopProtocol) DisconnectLink(link types.Link) error {
	return p.inner.DisconnectLink(link)
}

// ConnectSatellite forwards satellite connection to the inner protocol
func (p *IslAddSmartLoopProtocol) ConnectSatellite(s types.Node) error {
	return p.inner.ConnectSatellite(s)
}

// DisconnectSatellite forwards satellite disconnection to the inner protocol
func (p *IslAddSmartLoopProtocol) DisconnectSatellite(s types.Node) error {
	return p.inner.DisconnectSatellite(s)
}

// Links returns all candidate links from the inner protocol
func (p *IslAddSmartLoopProtocol) Links() []types.Link {
	return p.inner.Links()
}

// Established returns all currently established links
func (p *IslAddSmartLoopProtocol) Established() []types.Link {
	return p.inner.Established()
}

// UpdateLinks enhances the underlying protocol's MST with smart loops
func (p *IslAddSmartLoopProtocol) UpdateLinks() ([]types.Link, error) {
	if p.satellite == nil {
		return nil, errors.New("satellite not mounted")
	}

	p.mu.Lock()
	// Return cached if position hasn't changed
	if p.position.Equals(p.satellite.PositionVector()) {
		p.mu.Unlock()
		p.resetEvent.Wait() // Wait until ready
		return p.resultCache, nil
	}
	p.position = p.satellite.PositionVector()
	p.resetEvent.Reset() // Mark as busy
	p.mu.Unlock()

	innerEstablished, err := p.inner.UpdateLinks()
	if err != nil {
		return nil, err
	}

	satToLinks := make(map[types.Node]int)
	uniqueLinks := make(map[types.Node]types.Link)

	// Count how many times each satellite appears in links
	for _, link := range innerEstablished {
		n1, n2 := link.Nodes()
		for _, n := range []types.Node{n1, n2} {
			satToLinks[n]++
			if satToLinks[n] == 1 {
				uniqueLinks[n] = link
			} else {
				delete(uniqueLinks, n)
			}
		}
	}

	var s types.Node
	var additions []types.Link
	for _, candidate := range p.inner.Links() {
		if candidate.Distance() <= configs.MaxISLDistance {
			n1, n2 := candidate.Nodes()
			if uniqueLinks[n1] != nil {
				s = n1
			} else if uniqueLinks[n2] != nil {
				s = n2
			} else {
				continue
			}

			other := candidate.GetOther(s)
			if other == nil || satToLinks[s] >= p.config.Neighbours || satToLinks[other] >= p.config.Neighbours {
				continue
			}
			satToLinks[s]++
			satToLinks[other]++
			additions = append(additions, candidate)
		}
	}

	p.resultCache = append(innerEstablished, additions...)
	p.resetEvent.Set() // Mark as ready
	return p.resultCache, nil
}
